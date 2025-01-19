package mutation

import (
	"encoding/json"
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
	var expectedPatch interface{}
	var actualPatch interface{}
	err = json.Unmarshal([]byte(p), &expectedPatch)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(got, &actualPatch)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedPatch, actualPatch)
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
