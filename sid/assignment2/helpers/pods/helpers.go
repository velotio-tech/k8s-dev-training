package pods

import (
	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodRequestHandler(
	pod v1.Pod, client client.Client, log logr.Logger) (ctrl.Result, error) {
	log.Info("pod handler called with " + pod.Name)

	return ctrl.Result{}, nil
}

// Handles Creation of a new Pod
func HandlePodCreate(pod *v1.Pod, sanity *qav1.CodeSanity) {
	var watchPod bool
	for _, podName := range sanity.Spec.PodNames {
		if podName == pod.Name {
			watchPod = true
			break
		}
	}

	if !watchPod {
		return
	}

}
