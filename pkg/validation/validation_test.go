package validation

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidatePod(t *testing.T) {
	v := NewValidator(logger())

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
			Namespace: "some-namespace",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}

	val, err := v.ValidatePod(pod)
	assert.Nil(t, err)
	assert.True(t, val.Valid)

	pod2 := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
			Namespace: "protected",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}

	val2, err2 := v.ValidatePod(pod2)
	assert.Nil(t, err2)
	assert.False(t, val2.Valid)
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = io.Discard
	return mute.WithField("logger", "test")
}
