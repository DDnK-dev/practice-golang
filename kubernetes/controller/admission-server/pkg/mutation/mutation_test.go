package mutation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"webhook-server/pkg/mutation"
	"webhook-server/pkg/testutil"
)

func TestMutatePodPatch(t *testing.T) {
	m := mutation.NewMutator(testutil.Logger())
	got, err := m.MutatePodPatch(pod())
	if err != nil {
		t.Fatal(err)
	}
	p := patch()
	g := string(got)
	assert.Equal(t, g, p)
}

func patch() string {
	patch := `[
		{
			"value":
			[
				{
					"name": "KUBE",
					"value": "true"
				}
			],
			"op": "add",
			"path": "/spec/containers/0/env"
		},
		{
			"value":
			[
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "7"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "8"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "9"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "10"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "11"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "12"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "13"
				},
				{
					"effect": "NoSchedule",
					"key": "lifespan-remaining",
					"operator": "Equal",
					"value": "14"
				}
			],
			"op": "add",
			"path": "/spec/tolerations"
		}
	]`
	patch = strings.ReplaceAll(patch, "\n", "")
	patch = strings.ReplaceAll(patch, "\t", "")
	patch = strings.ReplaceAll(patch, " ", "")

	return patch
}

func pod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"lifespan-requested": "7",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "test",
				Image: "busybox",
			}},
		},
	}
}
