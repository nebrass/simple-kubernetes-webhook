package mutation

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestMutatePodPatch(t *testing.T) {
	m := NewMutator(testLogger())
	got, err := m.MutatePodPatch(pod())
	if err != nil {
		t.Fatal(err)
	}

	p := patch()
	g := string(got)
	assert.Equal(t, p, g)
}

func BenchmarkMutatePodPatch(b *testing.B) {
	m := NewMutator(testLogger())
	pod := pod()

	for i := 0; i < b.N; i++ {
		_, err := m.MutatePodPatch(pod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func pod() *corev1.Pod {
	return &corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}
}

func patch() string {
	patch := `[{"op":"replace","path":"/spec/containers/0/image","value":"contoso.acr.io/busybox"}]`

	return patch
}

func testLogger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = io.Discard
	return mute.WithField("logger", "test")
}
