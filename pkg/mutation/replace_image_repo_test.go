package mutation

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestImageMutationWithRepository(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "image-test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app", Image: "docker.io/library/busybox"},
			},
		},
	}

	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "image-test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app", Image: "contoso.acr.io/library/busybox"},
			},
		},
	}

	got, err := imageMutator{logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestImageMutationWithoutRepository(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "image-test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app", Image: "busybox"},
			},
		},
	}

	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "image-test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app", Image: "contoso.acr.io/busybox"},
			},
		},
	}

	got, err := imageMutator{logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestImageMutationIdempotence(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "image-test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app", Image: "contoso.acr.io/busybox"},
			},
		},
	}

	want := pod.DeepCopy()

	got, err := imageMutator{logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func logger() logrus.FieldLogger {
	return logrus.New()
}
