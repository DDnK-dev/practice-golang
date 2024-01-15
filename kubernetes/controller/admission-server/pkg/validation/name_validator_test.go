package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"webhook-server/pkg/testutil"
)

func TestValidation(t *testing.T) {
	nameValidator := nameValidator{
		Logger: testutil.Logger(),
	}
	tests := []struct {
		name string
		pod  *corev1.Pod
		want validation
	}{
		{
			name: "valid pod",
			pod:  pod("valid-pod"),
			want: validation{
				Valid:  true,
				Reason: "",
			},
		},
		{
			name: "invalid pod",
			pod:  pod("bad-name"),
			want: validation{
				Valid:  false,
				Reason: "pod name contains bad name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nameValidator.Validate(tt.pod)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, got)
		})

	}
}

func pod(name string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test-container",
					Image: "test-image",
				},
			},
		},
	}
}
