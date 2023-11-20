package admission

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webhook-server/pkg/mutation"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// admission 관련 로직을 구현을 위한 컨테이너
type Admitter struct {
	Logger  *logrus.Entry
	Request *admissionv1.AdmissionRequest
}

func (a Admitter) MutatePodReview() (*admissionv1.AdmissionReview, error) {
	pod, err := a.Pod()
	if err != nil {
		e := fmt.Sprintf("Failed to get pod from admission review request: %v", err)
		// uid types.UID, allowed bool, httpCode int32, reason string
		return reviewResponse(a.Request.UID, true, http.StatusBadRequest, e), err
	}
	// mutation 관련 로직은 mutation에 구현한다.
	m := mutation.NewMutator(a.Logger)
	patch, err := m.MutatePodPatch(pod)
	if err != nil {
		e := fmt.Sprintf("Failed to mutate pod: %v", err)
		return reviewResponse(a.Request.UID, true, http.StatusBadRequest, e), err
	}

	return patchReviewResponse(a.Request.UID, patch), nil
}

func (a Admitter) ValidatePodReview() (*admissionv1.AdmissionReview, error) {
	return nil, nil // TODO: implement this
}

// Pod extracts pod from request
func (a Admitter) Pod() (*corev1.Pod, error) {
	if a.Request.Kind.Kind != "Pod" {
		return nil, fmt.Errorf("expect resource to be Pod, got %s", a.Request.Kind.Kind)
	}
	var pod corev1.Pod
	if err := json.Unmarshal(a.Request.Object.Raw, &pod); err != nil {
		return nil, err
	}
	return &pod, nil
}

func reviewResponse(uid types.UID, allowed bool, httpCode int32, reason string) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
			Result: &metav1.Status{
				Code:    httpCode,
				Message: reason,
			},
		},
	}
}

func patchReviewResponse(uid types.UID, patch []byte) *admissionv1.AdmissionReview {
	return nil // TODO: implement this
}
