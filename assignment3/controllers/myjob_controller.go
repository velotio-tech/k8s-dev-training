/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	velotiov1 "github.com/pankaj9310/k8s-dev-training/pankaj/assignment3/api/v1"
	batch "k8s.io/api/batch/v1"
	batchbeta "k8s.io/api/batch/v1beta1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MyJobReconciler reconciles a MyJob object
type MyJobReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var (
	jobOwnerKey1 = ".metadata.controller.job"
	jobOwnerKey2 = ".metadata.controller.cronjob"
)

// +kubebuilder:rbac:groups=velotio.pankaj.io,resources=myjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=velotio.pankaj.io,resources=myjobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=velotiov1,resources=cronjobs,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=velotiov1,resources=jobs,verbs=get;list;watch;create;update;delete

func (r *MyJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("myjob", req.NamespacedName)

	log.Info("fetching MyJob resource")
	myJob := velotiov1.MyJob{}
	if err := r.Client.Get(ctx, req.NamespacedName, &myJob); err != nil {
		log.Error(err, "failed to get MyJob resource.")
		// Ignore NotFound errors as they will be retried automatically if the
		// resource is created in future.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if fmt.Sprintf("%s", myJob.Spec.ResourceType) == "job" {
		if err := r.cleanupOwnedResources(ctx, log, &myJob); err != nil {
			log.Error(err, "failed to clean up old Deployment resources for this MyJob.")
			return ctrl.Result{}, err
		}

	} else if fmt.Sprintf("%s", myJob.Spec.ResourceType) == "cronjob" {
		if err := r.cleanupOwnedCronjobResources(ctx, log, &myJob); err != nil {
			log.Error(err, "failed to clean up old Deployment resources for this Cronjob Myjob.")
			return ctrl.Result{}, err
		}
	}
	// log.Info("Job type", string(myJob.Spec.ResourceType))
	log = log.WithValues("Job_type", myJob.Spec.ResourceType)
	log.Info("checking if an existing Job exists for this resource")

	job := batch.Job{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: myJob.Namespace, Name: myJob.Spec.JobName}, &job)
	if apierrors.IsNotFound(err) {
		log.Info("could not find existing Job for MyJob, creating one...")
		if fmt.Sprintf("%s", myJob.Spec.ResourceType) == "job" {
			log.Info("Build Job")
			job = *buildJob(myJob)
			if err := r.Client.Create(ctx, &job); err != nil {
				log.Error(err, "failed to create Job resource")
				return ctrl.Result{}, err
			}
			r.Recorder.Eventf(&myJob, core.EventTypeNormal, "Created", "Created job %q", job.Name)
			log.Info("created Job resource for MyJob")
		}
		//process cronjob
		cronjob := batchbeta.CronJob{}
		err := r.Client.Get(ctx, client.ObjectKey{Namespace: myJob.Namespace, Name: myJob.Spec.JobName}, &cronjob)
		if apierrors.IsNotFound(err) {
			log.Info("could not find existing Job for MyJob, creating one...")

			if fmt.Sprintf("%s", myJob.Spec.ResourceType) == "cronjob" {
				cronjob := batchbeta.CronJob{}
				cronjob = *buildCronJob(myJob)
				log.Info("Build CronJob")
				if err := r.Client.Create(ctx, &cronjob); err != nil {
					log.Error(err, "failed to create CronJob resource")
					return ctrl.Result{}, err
				}
				r.Recorder.Eventf(&myJob, core.EventTypeNormal, "Created", "Created Cronjob %q", cronjob.Name)
				log.Info("created Job resource for MyJob")
			}
		}

		// r.Recorder.Eventf(&myJob, core.EventTypeNormal, "Created", "Created job %q", job.Name)
		// log.Info("created Job resource for MyJob")
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "failed to get Job for MyJob resource")
		return ctrl.Result{}, err
	}

	log.Info("existing Job resource already exists for MyJob, checking replica count")

	log.Info("resource status synced")

	return ctrl.Result{}, nil
}

func (r *MyJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addOwnerKey := func(object runtime.Object) []string {
		obj := object.(metav1.Object)
		owner := metav1.GetControllerOf(obj)
		if owner == nil {
			return nil
		}
		return []string{owner.Name}
	}

	if err := mgr.GetFieldIndexer().IndexField(&batch.Job{}, jobOwnerKey1, addOwnerKey); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(&batchbeta.CronJob{}, jobOwnerKey2, addOwnerKey); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&velotiov1.MyJob{}).
		Owns(&batchbeta.CronJob{}).
		Owns(&batch.Job{}).
		Complete(r)
}
