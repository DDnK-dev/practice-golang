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

type podValidator interface {
	Validate(*corev1.Pod) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// 각 validator에 대해 validate를 수행하고 결과 json을 반환한다
func (m *Validator) ValidatePodPatch(pod *corev1.Pod) (validation, error) {
	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else { // job을 만든다던가
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	m.Logger = m.Logger.WithField("pod_name", podName)
	m.Logger.Debug("start validation")

	validations := []podValidator{
		nameValidator{Logger: m.Logger},
	}

	for _, v := range validations {
		valid, err := v.Validate(pod)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !valid.Valid {
			return validation{Valid: false, Reason: valid.Reason}, nil
		}
	}
	return validation{Valid: true, Reason: ""}, nil
}
