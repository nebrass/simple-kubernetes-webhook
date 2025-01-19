package validation

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNamespaceValidatorProtectedNamespace(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "protected",
		},
	}

	validator := namespaceValidator{testLogger()}
	validation, err := validator.Validate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, validation.Valid)
	assert.Equal(t, "operation targeting namespace \"protected\" is not allowed", validation.Reason)
}

func TestNamespaceValidatorAllowedNamespace(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "allowed",
		},
	}

	validator := namespaceValidator{logger()}
	validation, err := validator.Validate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, validation.Valid)
	assert.Equal(t, "valid namespace", validation.Reason)
}

func testLogger() logrus.FieldLogger {
	return logrus.New()
}
