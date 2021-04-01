package controllers

import (
	"context"

	"github.com/go-logr/logr"
	velotiov1 "github.com/pankaj9310/k8s-dev-training/pankaj/assignment3/api/v1"
	batch "k8s.io/api/batch/v1"
	batchbeta "k8s.io/api/batch/v1beta1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// cleanupOwnedResources will Delete any existing Deployment resources that
// were created for the given MyJob that no longer match the
// myJob.spec.jobName field.
func (r *MyJobReconciler) cleanupOwnedResources(ctx context.Context, log logr.Logger, myJob *velotiov1.MyJob) error {
	log.Info("finding existing Deployments for MyJob Job resource")

	// List all job resources owned by this MyJob
	var jobs batch.JobList
	if err := r.List(ctx, &jobs, client.InNamespace(myJob.Namespace), client.MatchingFields{jobOwnerKey1: myJob.Name}); err != nil {
		log.Info("Error to get job list.")
		return err
	}

	deleted := 0
	for _, depl := range jobs.Items {
		if depl.Name == myJob.Spec.JobName {
			// If this job's name matches the one on the MyKind resource
			// then do not delete it.
			continue
		}

		if err := r.Client.Delete(ctx, &depl); err != nil {
			log.Error(err, "failed to delete Deployment resource")
			return err
		}

		r.Recorder.Eventf(myJob, core.EventTypeNormal, "Deleted", "Deleted job %q", depl.Name)
		deleted++
	}

	log.Info("finished cleaning up old Deployment resources", "number_deleted", deleted)

	return nil
}

// cleanupOwnedResources will Delete any existing Deployment resources that
// were created for the given MyJob that no longer match the
// myJob.spec.jobName field.
func (r *MyJobReconciler) cleanupOwnedCronjobResources(ctx context.Context, log logr.Logger, myJob *velotiov1.MyJob) error {
	log.Info("finding existing Deployments for MyKind Cronjob resource")
	// List all cronjobs resources owned by this MyJob
	var cronJobs batchbeta.CronJobList
	if err := r.List(ctx, &cronJobs, client.InNamespace(myJob.Namespace), client.MatchingFields{jobOwnerKey2: myJob.Name}); err != nil {
		log.Info("Error to get cronjob list")
		return err
	}

	deleted := 0
	for _, depl := range cronJobs.Items {
		if depl.Name == myJob.Spec.JobName {
			// If this job's name matches the one on the MyKind resource
			// then do not delete it.
			continue
		}

		if err := r.Client.Delete(ctx, &depl); err != nil {
			log.Error(err, "failed to delete Deployment resource")
			return err
		}

		r.Recorder.Eventf(myJob, core.EventTypeNormal, "Deleted", "Deleted job %q", depl.Name)
		deleted++
	}

	log.Info("finished cleaning up old Deployment resources", "number_deleted", deleted)

	return nil
}

func buildJob(myJob velotiov1.MyJob) *batch.Job {
	job := batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            myJob.Spec.JobName,
			Namespace:       myJob.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&myJob, velotiov1.GroupVersion.WithKind("MyJob"))},
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"example-controller.jetstack.io/job-name": myJob.Spec.JobName,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  "job-busybox",
							Image: "busybox",
							Command: []string{
								"/bin/sh",
								"-ec",
								"sleep 15",
							},
						},
					},
					RestartPolicy: core.RestartPolicyOnFailure,
				},
			},
		},
	}
	return &job
}

func buildCronJob(myJob velotiov1.MyJob) *batchbeta.CronJob {
	cronJob := batchbeta.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:            myJob.Spec.JobName,
			Namespace:       myJob.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&myJob, velotiov1.GroupVersion.WithKind("MyJob"))},
		},
		Spec: batchbeta.CronJobSpec{
			Schedule:          "*/1 * * * *",
			ConcurrencyPolicy: batchbeta.ForbidConcurrent,
			JobTemplate: batchbeta.JobTemplateSpec{
				Spec: batch.JobSpec{
					Template: core.PodTemplateSpec{
						Spec: core.PodSpec{
							RestartPolicy: core.RestartPolicyOnFailure,
							Containers: []core.Container{
								{
									Name:  "cronjob-busybox",
									Image: "busybox",
									Command: []string{
										"/bin/sh",
										"-ec",
										"sleep 15",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &cronJob
}
