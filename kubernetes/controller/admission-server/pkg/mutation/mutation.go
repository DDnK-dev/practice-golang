package mutation

// TODO: implement mutation logic

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
)

type Mutator struct {
	Logger *logrus.Entry
}

func NewMutator(logger *logrus.Entry) *Mutator {
	return &Mutator{
		Logger: logger,
	}
}

// podMutator is an interface used to group functions mutating pods
type podMutator interface {
	Mutate(*corev1.Pod) (*corev1.Pod, error)
	Name() string
}

// return json pod patch. Do all mutation method
func (m *Mutator) MutatePodPatch(pod *corev1.Pod) ([]byte, error) {
	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else { // job을 만든다던가 https://github.com/kubernetes/kubernetes/issues/44501
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)

	// list of all mutations to be applied to pod
	mutations := []podMutator{
		injectEnv{Logger: log},
		minLifespanTolerations{Logger: log},
	}

	mpod := pod.DeepCopy()

	// apply all mutations
	for _, m := range mutations {
		var err error
		mpod, err = m.Mutate(mpod)
		if err != nil {
			return nil, err
		}
	}

	// generate json patch
	patch, err := jsondiff.Compare(pod, mpod)
	if err != nil {
		return nil, err
	}
	patchb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}
	return patchb, nil
}
