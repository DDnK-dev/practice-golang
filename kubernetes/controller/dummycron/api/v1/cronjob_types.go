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

package v1

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConcurrencyPolicy describes how the job will be handled.
// Only one of the following concurrent policies may be specified.
// Ifnone of the following policies is specified, the default one is AllowConcurrent.
// +kubebuilder:validation:Enum=Allow;Forbid;Replace

type ConcurrencyPolicy string

const (
	// AllowConcurrent Cronjob이 동시에 실행되는 것을 허용한다
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent 동시에 실행되는 것을 허용하지 않는다. 만약 이전의 실행이 끝나지 않았다면 다음 실행을 건너뛴다
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent 현재 실행중인 Job을 취소하고 새로운 Job을 실행한다
	ReplaceConcurrent ConcurrencyPolicy = "Replace"
)

// CronJobSpec defines the desired state of CronJob
type CronJobSpec struct {
	//+kubebuilder:validation:MinLength=0
	// cron format으로 이루어져야 한다
	Schedule string `json:"schedule"`

	//+kubebuilder:validation:Minimum=0

	// Optional, 데드라인이 지나면 Job을 취소한다
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// 동시에 실행된 Job을 어떻게 처리할지 정한다
	// 허용되는 값은 다음과 같다.
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	// +optional
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// 이 플래그는 컨트럴러에게 다음 실행을 중지하라고 알려준다. 이미 시작된 실행에는 적용되지 않는다. 기본값은 false이다
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// CronJob이 실행될때 생성되는 Job의 템플릿이다
	JobTemplate batchv1.JobTemplateSpec `json:"jobTemplate"`

	//+kubebuilder:validation:Minimum=0

	// successful job을 몇개까지 유지할지 정한다. 이는 명시적인 0과 지정되지 않은 것을 구분하기 위해 포인터로 정의되어 있다
	// +optional
	SuccessfulJobsHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty"`

	//+kubebuilder:validation:Minimum=0

	// failed job을 몇개까지 유지할지 정한다. 이는 명시적인 0과 지정되지 않은 것을 구분하기 위해 포인터로 정의되어 있다
	// +optional
	FailedJobsHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

// CronJobStatus defines the observed state of CronJob
type CronJobStatus struct {
	// 현재 running중인 job에 대한 포인트 리스트
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// 마지막으로 job이 성공적으로 schedule 된 시간을 기록한다
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CronJob is the Schema for the cronjobs API
type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CronJobSpec   `json:"spec,omitempty"`
	Status CronJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CronJobList contains a list of CronJob
type CronJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CronJob{}, &CronJobList{})
}
