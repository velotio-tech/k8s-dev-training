package codesanity

import (
	"context"
	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
	"github.com/farkaskid/k8s-dev-training/assignment2/helpers/jobs"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
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

	log.Info("---Code Sanity Handler Called---")

	if err := checkTestJobs(ctx, c, sanity, log); err != nil {
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

	log.Info("---Code Sanity Handler Ended---")

	return ctrl.Result{}, nil
}

// Checks all the test jobs that were spawned either directly by CodeSanity CR or by cronjobs spawned by
// CodeSanity CR. For each job, it checks if the job was completed and whether the tests were run successfully.
// It adds the pod corresponding to the test job to healthy or unhealthy status lists accordingly.
func checkTestJobs(ctx context.Context, c client.Client, sanity *qav1.CodeSanity, log logr.Logger) error {

	log.Info("<-Checking spawned jobs->")

	var completedJobs []batchv1.Job

	testingJobs, err := getTestJobs(ctx, c, sanity, log)
	if err != nil {
		return err
	}

	// Checking completed jobs for test results
	for _, job := range testingJobs {
		finished, condition := jobs.IsJobFinished(&job)

		if !finished {
			log.Info("Skipping running job: " + job.Name)
			continue
		}

		completedJobs = append(completedJobs, job)
		testedPod := job.Labels["pod"]

		switch condition {
		case batchv1.JobFailed:
			log.Info("Test for pod: " + testedPod + " has failed")
			sanity.AddUnhealthyPod(testedPod)

		case batchv1.JobComplete:
			log.Info("Test for pod: " + testedPod + " has passed")
			sanity.AddHealthyPod(testedPod)
		}
	}

	// Removing completed jobs if jobs were spawned by code sanity
	if sanity.Spec.Resource == qav1.Job {
		for _, job := range completedJobs {
			if err := c.Delete(ctx, &job, client.PropagationPolicy(metav1.DeletePropagationBackground)); err != nil {
				log.Error(err, "Failed to delete completed job: "+job.Name)

				return err
			}
		}
	}

	return nil
}

// Gets all the jobs that were spawned either directly by CodeSanity CR or indirectly by the cronjobs created by
// the CodeSanity CR
func getTestJobs(
	ctx context.Context, c client.Client, sanity *qav1.CodeSanity, log logr.Logger) ([]batchv1.Job, error) {
	log.Info("<-Getting all spawned jobs->")

	var testJobs []batchv1.Job

	// Get jobs that were directly spawned by CodeSanity CR
	if sanity.Spec.Resource == qav1.Job {
		log.Info("getting jobs for: " + sanity.Name)

		var testJobList batchv1.JobList

		if err := c.List(
			ctx,
			&testJobList,
			client.InNamespace(sanity.Namespace),
			client.MatchingFields{JobOwnerKey: sanity.Name},
		); err != nil {
			log.Error(err, "Failed to get jobs for: "+sanity.Name)

			return testJobs, err
		}

		return testJobList.Items, nil
	}

	// Get the cronjobs spawned by CodeSanity CR and then get the jobs spawned by those cronjobs
	var testCronJobList v1beta1.CronJobList

	log.Info("getting cronjobs for: " + sanity.Name)

	if err := c.List(
		ctx,
		&testCronJobList,
		client.InNamespace(sanity.Namespace),
		client.MatchingFields{JobOwnerKey: sanity.Name},
	); err != nil {
		log.Error(err, "Failed to get cronjobs for: "+sanity.Name)

		return testJobs, err
	}

	for _, cronjob := range testCronJobList.Items {
		var testJobList batchv1.JobList

		log.Info("getting jobs for: " + cronjob.Name)

		if err := c.List(
			ctx,
			&testJobList,
			client.InNamespace(sanity.Namespace),
			client.MatchingFields{JobOwnerKey: cronjob.Name},
		); err != nil {
			log.Error(err, "Failed to get jobs for: "+cronjob.Name)

			return testJobs, err
		}

		testJobs = append(testJobs, testJobList.Items...)
	}

	return testJobs, nil
}

// Spawns new jobs/cronjobs for the pods selected according to different attributes in CodeSanity CR spec
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
