package validation

// TODO: implement validation logic

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type Validator struct {
	Logger *logrus.Entry
}

func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{
		Logger: logger,
	}
}

func (m *Validator) ValidatePodPatch(pod *corev1.Pod) ([]byte, error) {
	return []byte{}, nil // do nothing now
}
