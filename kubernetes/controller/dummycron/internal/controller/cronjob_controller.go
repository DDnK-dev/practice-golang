/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron"
	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "my.kubebuilder/dummycron/api/v1"
)

/*
이해를 돕기 위해 임의의 cronjob status를 첨부한다

- active job 있을 때 status
status:
  active:
  - apiVersion: batch/v1
    kind: Job
    name: hello-28321329
    namespace: default
    resourceVersion: "2594"
    uid: e6e77da6-3945-43b8-8216-bcbefdc085cd
  lastScheduleTime: "2023-11-06T14:09:00Z"
  lastSuccessfulTime: "2023-11-06T14:08:03Z"

- idle 일떄의 status
status:
  lastScheduleTime: "2023-11-06T14:05:00Z"
  lastSuccessfulTime: "2023-11-06T14:05:08Z"
*/

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Clock
}

// 대충 time.Now를 부르는 clock을 만들 것
type realClock struct{}

func (_ realClock) Now() time.Time { return time.Now() }

type Clock interface {
	Now() time.Time
}

// 몇가지 RBAC이 더 필요하다 (Job을 컨트롤하기 위해)
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch

//+kubebuilder:rbac:groups=batch.my.kubebuilder,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.my.kubebuilder,resources=cronjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.my.kubebuilder,resources=cronjobs/finalizers,verbs=update

// 필요한 annotation을 정의하자
var (
	scheduledTimeAnnotation = "batch.tutorial.kubebuilder.io/scheduled-at"
)

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CronJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1. CronJob 이름 가져오기
	var cronJob batchv1.CronJob
	if err := r.Get(ctx, req.NamespacedName, &cronJob); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// not-found error는 무시하자 (즉각적인 requeue로는 처리가 어렵고, 다음 reconcile에서 처리하도록 하자)
		// 그리고 delete request에서는 오브젝트를 받을 수 있을 것
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. 모든 active job List를 가져오고, status 업데이트하기
	var childJobs kbatch.JobList
	// MatchFields를 사용하여, owner가 현재 cronjob인 job을 가져온다 (퍼포먼스를 위해 jobOwnerKey를 인덱스로 사용할것. 아래 Setup부분을 참고)
	if err := r.List(ctx, &childJobs, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}

	// active job 찾기
	var activeJobs []*kbatch.Job
	var successfulJobs []*kbatch.Job
	var failedJobs []*kbatch.Job
	var mostRecentTime *time.Time // 마지막에 실행된 job의 시간을 저장할 것

	// is Job finished?
	isJobFinished := func(job *kbatch.Job) (bool, kbatch.JobConditionType) {
		for _, c := range job.Status.Conditions { // status의 condition 필드를 순회하며 조건 확인
			if (c.Type == kbatch.JobComplete || c.Type == kbatch.JobFailed) && c.Status == corev1.ConditionTrue {
				return true, c.Type
			}
		}
		return false, ""
	}

	// get ScheduledTime for job
	getScheduledTimeForJob := func(job *kbatch.Job) (*time.Time, error) {
		timeRaw := job.Annotations[scheduledTimeAnnotation]
		if timeRaw == "" {
			return nil, nil
		}
		timeParsed, err := time.Parse(time.RFC3339, timeRaw)
		if err != nil {
			return nil, err
		}
		return &timeParsed, nil
	}

	//
	for i := range childJobs.Items {
		_, finishedType := isJobFinished(&childJobs.Items[i])
		switch finishedType {
		case "": // not finished
			activeJobs = append(activeJobs, &childJobs.Items[i])
		case kbatch.JobComplete:
			successfulJobs = append(successfulJobs, &childJobs.Items[i])
		case kbatch.JobFailed:
			failedJobs = append(failedJobs, &childJobs.Items[i])
		}

		// annotation 상의 launch time을 가져와서, 다음 job을 실행할 시간을 계산한다
		scheduledTimeForJob, err := getScheduledTimeForJob(&childJobs.Items[i])
		if err != nil {
			log.Error(err, "unable to parse schedule time for child job", "job", &childJobs.Items[i])
			continue
		}
		if scheduledTimeForJob != nil {
			if mostRecentTime == nil || mostRecentTime.Before(*scheduledTimeForJob) {
				mostRecentTime = scheduledTimeForJob
			}
		}
	}

	if mostRecentTime != nil {
		cronJob.Status.LastScheduleTime = &metav1.Time{Time: *mostRecentTime} // CronJob의 Status를 업데이트
	} else {
		cronJob.Status.LastScheduleTime = nil
	}

	cronJob.Status.Active = nil
	for _, active := range activeJobs {
		jobRef, err := ref.GetReference(r.Scheme, active)
		if err != nil {
			log.Error(err, "unable to make reference to active job", "job", active)
			continue
		}
		cronJob.Status.Active = append(cronJob.Status.Active, *jobRef)
	}
	log.V(1).Info("job count", "active jobs", len(activeJobs), "successful jobs", len(successfulJobs), "failed jobs", len(failedJobs))

	// 상태 업데이트
	if err := r.Status().Update(ctx, &cronJob); err != nil {
		log.Error(err, "unable to update CronJob status")
		return ctrl.Result{}, err
	}

	// 3. history limit에 따라 old job 삭제하기
	// NB: best effort basis로, 실패해도 무시하자
	if cronJob.Spec.FailedJobsHistoryLimit != nil {
		sort.Slice(failedJobs, func(i, j int) bool { // failedJobs를 실행 시간 순으로 정렬
			if failedJobs[i].Status.StartTime == nil {
				return failedJobs[j].Status.StartTime != nil
			}
			return failedJobs[i].Status.StartTime.Before(failedJobs[j].Status.StartTime)
		})
		for i, job := range failedJobs {
			if int32(i) >= int32(len(failedJobs))-*cronJob.Spec.FailedJobsHistoryLimit {
				break
			}
			if err := r.Delete(ctx, job, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				log.Error(err, "unable to delete old failed job", "job", job)
			} else {
				log.V(0).Info("deleted old failed ")
			}
		}
	}

	// 4. check suspended
	if cronJob.Spec.Suspend != nil && *cronJob.Spec.Suspend {
		log.V(1).Info("cronjob suspended, skipping")
		return ctrl.Result{}, nil
	}

	// 5. get next scheduled run
	// cron library를 사용하여 다음 스케줄 시간을 계산한다. 마지막 실행시간, 혹은 CronJob 실행 시간으로부터 얼마나 지났나?
	getNextSchedule := func(cronJob *batchv1.CronJob, now time.Time) (lastMissed time.Time, scheduled time.Time, err error) {
		sched, err := cron.ParseStandard(cronJob.Spec.Schedule)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		// 최적화용으로 그냥...
		var earliestTime time.Time
		if cronJob.Status.LastScheduleTime != nil {
			earliestTime = cronJob.Status.LastScheduleTime.Time
		} else {
			earliestTime = cronJob.ObjectMeta.CreationTimestamp.Time
		}
		if cronJob.Spec.StartingDeadlineSeconds != nil {
			// deadline이 있으면, deadline을 넘어서는 스케줄은 무시한다
			schedulingDeadline := now.Add(-time.Second * time.Duration(*cronJob.Spec.StartingDeadlineSeconds))
			if schedulingDeadline.After(earliestTime) {
				earliestTime = schedulingDeadline
			}
		}
		if earliestTime.After(now) {
			return time.Time{}, sched.Next(now), nil
		}

		starts := 0 // 시작 횟수
		for t := sched.Next(earliestTime); !t.After(now); t = sched.Next(t) {
			lastMissed = t
			starts++
			if starts > 100 { // 너무 오랫동안 job이 밀린 경우에 대비
				return time.Time{}, time.Time{}, fmt.Errorf("Too many missed start times (> 100). Set or decrease .spec.startingDeadlineSeconds or check clock skew")
			}
		}
		return lastMissed, sched.Next(now), nil
	}
	missedRun, nextRun, err := getNextSchedule(&cronJob, r.Now())
	if err != nil {
		log.Error(err, "unable to figure out next run time")
		return ctrl.Result{}, nil
	}
	scheduleResult := ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())} // save this so we can re-use it later
	log = log.WithValues("now", r.Now(), "next run", nextRun)

	// 6. schedule 에 job이 있다면 실행하기 (deadline이 지나지 않았고, concurrency policy에 저촉되지 않는다면)
	if missedRun.IsZero() {
		log.V(1).Info("no upcoming scheduled tims, sleeping until next run", "next", nextRun)
		return scheduleResult, nil
	}
	// make sure we're not to late to start the run (deadline)
	log = log.WithValues("current run", missedRun)
	tooLate := false
	if cronJob.Spec.StartingDeadlineSeconds != nil {
		tooLate = missedRun.Add(time.Duration(*cronJob.Spec.StartingDeadlineSeconds) * time.Second).Before(r.Now())
	}
	if tooLate {
		log.V(1).Info("missed starting deadline for last run, sleeping til next")
		return scheduleResult, nil
	}
	// 이제 concurrency policy를 확인해보자
	if
	// set up a real clock, since we're not in a test

	cronJob.Spec.ConcurrencyPolicy == batchv1.ForbidConcurrent && len(activeJobs) > 0 {
		log.V(1).Info("concurrency policy forbid concurrent, skipping")
		return scheduleResult, nil
	}

	if cronJob.Spec.ConcurrencyPolicy == batchv1.ReplaceConcurrent {
		for _, activeJob := range activeJobs {
			// we don't care if the job was already deleted
			if err := r.Delete(ctx, activeJob, client.PropagationPolicy(metav1.DeletePropagationBackground)); client.IgnoreNotFound(err) != nil {
				log.Error(err, "unable to delete old job", "job", activeJob)
				return ctrl.Result{}, err
			}
		}
	}
	// 걸리는게 없으면 job을생성하자
	// set up a real clock, since we're not in a test

	constructJobForCronJob := func(cronJob *batchv1.CronJob, scheduledTime time.Time) (*kbatch.Job, error) {
		name := fmt.Sprintf("%s-%d", cronJob.Name, scheduledTime.Unix())
		job := &kbatch.Job{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      make(map[string]string),
				Annotations: make(map[string]string),
				Name:        name,
				Namespace:   cronJob.Namespace,
			},
			Spec: *cronJob.Spec.JobTemplate.Spec.DeepCopy(),
		}
		for k, v := range cronJob.Spec.JobTemplate.Annotations {
			job.Annotations[k] = v
		}
		for k, v := range cronJob.Spec.JobTemplate.Labels {
			job.Labels[k] = v
		}
		if err := ctrl.SetControllerReference(cronJob, job, r.Scheme); err != nil {
			return nil, err
		}
		return job, nil
	}

	job, err := constructJobForCronJob(&cronJob, missedRun)
	if err != nil {
		log.Error(err, "unable to construct job from template")
		return scheduleResult, nil
	}

	// 실제 job 생성
	if err := r.Create(ctx, job); err != nil {
		log.Error(err, "unable to create Job for CronJob", "job", job)
		return ctrl.Result{}, err
	}
	log.V(1).Info("created Job for CronJob run", "job", job)

	// 다음 실행시간을 위하여 requeue
	return scheduleResult, nil
}

var (
	jobOwnerKey = ".metadata.controller" // job의 owner를 찾기 위한 인덱스
	apiGVStr    = batchv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// set up a real clock, since we're not in a test
	if r.Clock == nil {
		r.Clock = realClock{}
	}

	// 이 부분 잘 모르겠는데??????? field idnexer 쪽에서 저장할때 필터 역할을 하는 것인가??
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &kbatch.Job{}, jobOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner...
		job := rawObj.(*kbatch.Job)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}
		// ...make sure it's a CronJob...
		if owner.APIVersion != apiGVStr || owner.Kind != "CronJob" {
			return nil
		}
		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CronJob{}).
		Owns(&kbatch.Job{}).
		Complete(r)
}
