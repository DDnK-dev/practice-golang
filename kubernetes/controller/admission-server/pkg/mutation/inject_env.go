package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// pod에 환경변수를 주입하기 위한 mutator
type injectEnv struct {
	Logger logrus.FieldLogger
}

var _ podMutator = (*injectEnv)(nil) // interface 구현을 검증하기 위함

func (se injectEnv) Name() string {
	return "inject_env"
}

func (se injectEnv) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", se.Name())
	mpod := pod.DeepCopy()

	// build out env var slice
	envVars := []corev1.EnvVar{{
		Name:  "KUBE",
		Value: "true",
	}}

	// inject env vars into pod
	for _, envVar := range envVars {
		se.Logger.Debugf("pod env injected %s", envVar)
		injectEnvVar(mpod, envVar)
	}
	return mpod, nil
}

func injectEnvVar(pod *corev1.Pod, envVar corev1.EnvVar) {
	for i, container := range pod.Spec.Containers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.Containers[i].Env = append(container.Env, envVar)
		}
	}
	for i, container := range pod.Spec.InitContainers {
		if !HasEnvVar(container, envVar) {
			pod.Spec.InitContainers[i].Env = append(container.Env, envVar)
		}
	}
}

func HasEnvVar(container corev1.Container, envVar corev1.EnvVar) bool {
	for _, env := range container.Env {
		if env.Name == envVar.Name {
			return true
		}
	}
	return false
}
