package codesanity

import (
	"context"
	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
	"github.com/farkaskid/k8s-dev-training/assignment2/helpers/jobs"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const JobOwnerKey string = ".metadata.controller"

func CodeSanityRequestHandler(
	ctx context.Context, sanity *qav1.CodeSanity, c client.Client, scheme *runtime.Scheme, log logr.Logger) (ctrl.Result, error) {

	log.Info("---Code Sanity Handler Called----")

	if err := checkTestJobs(ctx, c, scheme, sanity, log); err != nil {
		return ctrl.Result{}, err
	}

	if err := spawnNewTestJobs(ctx, c, scheme, sanity, log); err != nil {
		return ctrl.Result{}, err
	}

	sanity.Status.LastRunAt = metav1.NewTime(time.Now())

	if err := c.Status().Update(ctx, sanity); err != nil {
		log.Error(err, "Failed to update sanity's status")

		return ctrl.Result{}, err
	}

	log.Info("---Code Sanity Handler Ended----")

	return ctrl.Result{}, nil
}

func checkTestJobs(
	ctx context.Context, c client.Client, scheme *runtime.Scheme, sanity *qav1.CodeSanity, log logr.Logger) error {

	log.Info("<-Checking spawned jobs->")

	var testJobs batchv1.JobList
	if err := c.List(
		ctx,
		&testJobs,
		client.InNamespace(sanity.Namespace),
		client.MatchingFields{JobOwnerKey: sanity.Name},
	); err != nil {
		return err
	}

	for _, job := range testJobs.Items {
		finished, condition := jobs.IsJobFinished(&job)

		if !finished {
			log.Info("Skipping running job: " + job.Name)
			continue
		}

		switch condition {
		case batchv1.JobFailed:
			log.Info(job.Name + " has failed")
		case batchv1.JobComplete:
			log.Info(job.Name + " has success")
		}
	}

	return nil
}

func spawnNewTestJobs(
	ctx context.Context, c client.Client, scheme *runtime.Scheme, sanity *qav1.CodeSanity, log logr.Logger) error {

	log.Info("<-Spawning new jobs->")

	var podList v1.PodList
	if err := c.List(ctx, &podList, client.InNamespace(sanity.Namespace)); err != nil {
		log.Error(err, "Failed to get podlist")

		return err
	}

	for _, pod := range podList.Items {
		if !sanity.PodValid(pod) {
			log.Info("Skipping pod: " + pod.Name)
			continue
		}

		if err := jobs.CreateJobForPod(ctx, &pod, c, scheme, sanity, log); err != nil {
			log.Error(err, "Failed to create job")

			return err
		}

		sanity.Status.ProcessedPods = append(sanity.Status.ProcessedPods, pod.Name)
	}
	return nil
}
