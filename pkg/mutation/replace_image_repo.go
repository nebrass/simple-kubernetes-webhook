package mutation

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// imageMutator is a container for image mutation logic
type imageMutator struct {
    Logger logrus.FieldLogger
}

// imageMutator implements the podMutator interface
var _ podMutator = (*imageMutator)(nil)

// Name returns the imageMutator short name
func (im imageMutator) Name() string {
    return "image_mutator"
}

// Mutate mutates the pod's container images to ensure they use the specified repository
func (im imageMutator) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
    const repositoryURL = "contoso.acr.io"

    im.Logger = im.Logger.WithField("mutation", im.Name())
    mpod := pod.DeepCopy()

    for i, container := range mpod.Spec.Containers {
        im.Logger.WithField("container", container.Name).Printf("Processing container image: %s", container.Image)

        // Split the image into repository and name:tag
        parts := strings.SplitN(container.Image, "/", 2)
        
        if len(parts) == 2 && strings.Contains(parts[0], ".") {
            // Image has a repository, replace it
            container.Image = fmt.Sprintf("%s/%s", repositoryURL, parts[1])
            im.Logger.WithField("container", container.Name).Printf("Replacing repository, new image: %s", container.Image)
        } else {
            // Image does not have a repository, prepend the default one
            container.Image = fmt.Sprintf("%s/%s", repositoryURL, container.Image)
            im.Logger.WithField("container", container.Name).Printf("Appending repository, new image: %s", container.Image)
        }

        mpod.Spec.Containers[i].Image = container.Image
    }

    return mpod, nil
}