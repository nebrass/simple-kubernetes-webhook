package validation

import (
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// namespaceValidator is a container for validating the namespace of pods
type namespaceValidator struct {
	Logger logrus.FieldLogger
}

// namespaceValidator implements the podValidator interface
var _ podValidator = (*namespaceValidator)(nil)

// Name returns the name of namespaceValidator
func (nv namespaceValidator) Name() string {
	return "namespace_validator"
}

// Validate inspects the namespace of a given pod and returns validation.
// The returned validation is invalid if the pod is targeting the "protected" namespace.
func (nv namespaceValidator) Validate(pod *corev1.Pod) (validation, error) {
	protectedNamespace := "protected"

	if pod.Namespace == protectedNamespace {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("operation targeting namespace %q is not allowed", protectedNamespace),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid namespace"}, nil
}
