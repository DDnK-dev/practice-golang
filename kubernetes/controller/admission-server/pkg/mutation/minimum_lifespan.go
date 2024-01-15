package mutation

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type minLifespanTolerations struct {
	Logger logrus.FieldLogger
}

var _ podMutator = (*minLifespanTolerations)(nil) // interface 구현을 검증하기 위함

func (mpl minLifespanTolerations) Name() string {
	return "min_lifespan_tolerations"
}

func (mpl minLifespanTolerations) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	const (
		lifespanLabel = "lifespan-requested"
		taintKey      = "lifespan-remaining"
		taintMaxAge   = 14
	)
	mpl.Logger = mpl.Logger.WithField("mutation", mpl.Name())
	mpod := pod.DeepCopy()

	if pod.Labels == nil || pod.Labels[lifespanLabel] == "" {
		mpl.Logger.WithField("min_lifespan", 0).
			Printf("no lifespan label foudn, applying default lifespan toleration")
		tn := []corev1.Toleration{{
			Key:      taintKey,
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		}}

		mpod.Spec.Tolerations = appendToleration(tn, mpod.Spec.Tolerations)
		return mpod, nil
	}
	ts := pod.Labels[lifespanLabel]
	minAge, err := strconv.Atoi(ts)
	if err != nil {
		return nil, fmt.Errorf("pod lifespan label %q is not an integer: %v", ts, err)
	}
	mpl.Logger.WithField("min_lifespan", ts).Printf("setting lifespan toleration")

	t := []corev1.Toleration{}
	for i := taintMaxAge; i >= minAge; i-- {
		t = appendToleration(t, []corev1.Toleration{{
			Key:      taintKey,
			Operator: corev1.TolerationOpEqual,
			Value:    strconv.Itoa(i),
			Effect:   corev1.TaintEffectNoSchedule,
		}})
	}

	mpod.Spec.Tolerations = appendToleration(t, mpod.Spec.Tolerations)
	return mpod, nil
}

func appendToleration(new, existing []corev1.Toleration) []corev1.Toleration {
	var toAppend []corev1.Toleration
	for _, n := range new {
		found := false
		for _, e := range existing {
			if reflect.DeepEqual(n, e) {
				found = true
				break
			}
		}
		if !found {
			toAppend = append(toAppend, n)
		}
	}
	return append(existing, toAppend...)
}
