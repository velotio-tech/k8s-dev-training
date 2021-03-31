package codesanity

import (
	"context"
	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
	"github.com/farkaskid/k8s-dev-training/assignment2/helpers/jobs"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func CodeSanityRequestHandler(
	ctx context.Context, sanity qav1.CodeSanity, c client.Client, scheme *runtime.Scheme, log logr.Logger) (ctrl.Result, error) {

	log.Info("codesanity handler called with " + sanity.Name)

	var podList v1.PodList
	if err := c.List(ctx, &podList, client.InNamespace(sanity.Namespace)); err != nil {
		log.Error(err, "Failed to get podlist")

		return ctrl.Result{}, err
	}

	for _, pod := range podList.Items {
		if !sanity.PodValid(pod) {
			log.Info("Skipping pod: " + pod.Name)
			continue
		}

		if err := jobs.CreateJobForPod(ctx, &pod, c, scheme, sanity, log); err != nil {
			log.Error(err, "Failed to create job")

			return ctrl.Result{}, err
		}

		sanity.Status.ProcessedPods = append(sanity.Status.ProcessedPods, pod.Name)
	}

	sanity.Status.LastRunAt = metav1.NewTime(time.Now())

	if err := c.Status().Update(ctx, &sanity); err != nil {
		log.Error(err, "Failed to update sanity's status")

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
