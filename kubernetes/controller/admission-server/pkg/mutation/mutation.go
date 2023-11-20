package mutation

// TODO: implement mutation logic

import (
	"github.com/sirupsen/logrus"
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

func (m *Mutator) MutatePodPatch(pod *corev1.Pod) ([]byte, error) {
	return []byte{}, nil // do nothing now
}
