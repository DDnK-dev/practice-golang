package validation

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

var _ podValidator = (*nameValidator)(nil)

type nameValidator struct {
	Logger logrus.FieldLogger
}

func (nv nameValidator) Name() string {
	return "name_validator"
}

func (nv nameValidator) Validate(pod *corev1.Pod) (validation, error) {
	badNames := []string{"bad-name", "bad-name2"}

	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else if pod.ObjectMeta.GenerateName != "" {
		podName = pod.ObjectMeta.GenerateName
	} else {
		return validation{}, fmt.Errorf("pod name is empty")
	}

	for _, badName := range badNames {
		if strings.Contains(podName, badName) {
			return validation{
				Valid:  false,
				Reason: "pod name contains bad name",
			}, nil
		}
	}
	return validation{
		Valid:  true,
		Reason: "",
	}, nil
}
