/*
Copyright 2026.

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

package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	assignmentsv1 "k8s-dev-training.com/assignments/api/v1"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	managedNamespace         = "default"
	workloadLabel            = "app"
	workloadValue            = "resourcemanager-workload"
	resourceManagerLabel     = "k8s-dev-training/resource-manager"
	jobIndexField            = ".metadata.resourceManagerName"
	resourceManagerFinalizer = "k8s-dev-training/job-cleanup"
)

// ResourceManagerReconciler reconciles a ResourceManager object
type ResourceManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=assignments.k8s-dev-training.com,resources=resourcemanagers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=assignments.k8s-dev-training.com,resources=resourcemanagers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=assignments.k8s-dev-training.com,resources=resourcemanagers/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

func (r *ResourceManagerReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Defense in depth: do not reconcile other namespaces.
	if req.Namespace != managedNamespace {
		return ctrl.Result{}, nil
	}

	var rm assignmentsv1.ResourceManager
	if err := r.Get(ctx, req.NamespacedName, &rm); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !rm.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&rm, resourceManagerFinalizer) {
			return ctrl.Result{}, nil
		}

		return r.finalizeResourceManager(ctx, &rm)
	}

	// Add the finalizer before creating resources that require cleanup.
	if !controllerutil.ContainsFinalizer(&rm, resourceManagerFinalizer) {
		base := rm.DeepCopy()
		controllerutil.AddFinalizer(&rm, resourceManagerFinalizer)

		if err := r.Patch(ctx, &rm, client.MergeFrom(base)); err != nil {
			return ctrl.Result{}, fmt.Errorf("add ResourceManager finalizer: %w", err)
		}

		// Reconcile again using the persisted object.
		return ctrl.Result{Requeue: true}, nil
	}

	// Discover associated Jobs during normal reconciliation too.
	jobs, err := r.listAssociatedJobs(ctx, &rm)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info(
		"Discovered associated Jobs",
		"resourceManager", rm.Name,
		"count", len(jobs),
	)

	deploymentName := rm.Name + "-workload"
	key := types.NamespacedName{
		Namespace: rm.Namespace,
		Name:      deploymentName,
	}

	var existing appsv1.Deployment
	err = r.Get(ctx, key, &existing)

	if apierrors.IsNotFound(err) {
		desired := desiredDeployment(&rm, deploymentName)

		if err := controllerutil.SetControllerReference(
			&rm, desired, r.Scheme,
		); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, desired); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("Created Deployment", "deployment", deploymentName)
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	desired := desiredDeployment(&rm, deploymentName)

	// Update only fields this controller owns.
	changed := false

	if !reflect.DeepEqual(existing.Labels, desired.Labels) {
		existing.Labels = desired.Labels
		changed = true
	}

	if !reflect.DeepEqual(existing.Spec.Replicas, desired.Spec.Replicas) {
		existing.Spec.Replicas = desired.Spec.Replicas
		changed = true
	}

	if !reflect.DeepEqual(existing.Spec.Selector, desired.Spec.Selector) {
		existing.Spec.Selector = desired.Spec.Selector
		changed = true
	}

	if !reflect.DeepEqual(existing.Spec.Template.Labels, desired.Spec.Template.Labels) {
		existing.Spec.Template.Labels = desired.Spec.Template.Labels
		changed = true
	}

	if !reflect.DeepEqual(existing.Spec.Template.Spec.Containers, desired.Spec.Template.Spec.Containers) {
		existing.Spec.Template.Spec.Containers = desired.Spec.Template.Spec.Containers
		changed = true
	}

	if changed {
		if err := r.Update(ctx, &existing); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Updated Deployment", "deployment", deploymentName)
	}

	return ctrl.Result{}, nil
}

func desiredDeployment(
	rm *assignmentsv1.ResourceManager,
	name string,
) *appsv1.Deployment {
	replicas := int32(1)

	labels := map[string]string{
		workloadLabel: workloadValue,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: rm.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "job-creator",
					Containers: []corev1.Container{
						{
							Name:  "job-helper",
							Image: "job-helper:dev",
							Env: []corev1.EnvVar{
								{
									Name:  "RESOURCE_MANAGER_NAME",
									Value: rm.Name,
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "JOB_COMMAND",
									Value: mustJSON(rm.Spec.JobCommand),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ResourceManagerReconciler) finalizeResourceManager(
	ctx context.Context,
	rm *assignmentsv1.ResourceManager,
) (ctrl.Result, error) {
	jobs, err := r.listAssociatedJobs(ctx, rm)
	if err != nil {
		// Keep the finalizer if cleanup cannot be verified.
		return ctrl.Result{}, err
	}

	if len(jobs) > 0 {
		for i := range jobs {
			job := &jobs[i]

			err := r.Delete(ctx, job)
			if err != nil && !apierrors.IsNotFound(err) {
				return ctrl.Result{}, fmt.Errorf(
					"delete associated Job %s/%s: %w",
					job.Namespace,
					job.Name,
					err,
				)
			}
		}

		// Wait for deletion to be reflected in the cache. The finalizer
		// stays until a subsequent reconciliation sees no associated Jobs.
		return ctrl.Result{RequeueAfter: 2 * time.Second}, nil
	}

	base := rm.DeepCopy()
	controllerutil.RemoveFinalizer(rm, resourceManagerFinalizer)

	if err := r.Patch(ctx, rm, client.MergeFrom(base)); err != nil {
		return ctrl.Result{}, fmt.Errorf("remove ResourceManager finalizer: %w", err)
	}

	return ctrl.Result{}, nil
}

// listAssociatedJobs finds Jobs using the cache index registered in
// SetupWithManager. It does not rely on owner references.
func (r *ResourceManagerReconciler) listAssociatedJobs(
	ctx context.Context,
	rm *assignmentsv1.ResourceManager,
) ([]batchv1.Job, error) {
	var jobList batchv1.JobList

	err := r.List(
		ctx,
		&jobList,
		client.InNamespace(rm.Namespace),
		client.MatchingFields{
			jobIndexField: rm.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("list associated Jobs: %w", err)
	}

	var associated []batchv1.Job

	for i := range jobList.Items {
		job := &jobList.Items[i]

		// This design expects these Jobs to have no owner references.
		// Never delete a Job that may belong to another owner.
		if len(job.OwnerReferences) != 0 {
			return nil, fmt.Errorf(
				"Job %s/%s has owner references; refusing to manage it",
				job.Namespace,
				job.Name,
			)
		}

		associated = append(associated, *job)
	}

	return associated, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&batchv1.Job{},
		jobIndexField,
		func(obj client.Object) []string {
			value := obj.GetLabels()[resourceManagerLabel]
			if value == "" {
				return nil
			}
			return []string{value}
		},
	); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&assignmentsv1.ResourceManager{}).
		Owns(&appsv1.Deployment{}).
		Named("resourcemanager").
		Watches(
			&batchv1.Job{},
			handler.EnqueueRequestsFromMapFunc(
				func(ctx context.Context, obj client.Object) []reconcile.Request {
					if obj.GetNamespace() != managedNamespace {
						return nil
					}

					name := obj.GetLabels()[resourceManagerLabel]
					if name == "" {
						return nil
					}

					return []reconcile.Request{
						{
							NamespacedName: types.NamespacedName{
								Namespace: obj.GetNamespace(),
								Name:      name,
							},
						},
					}
				},
			),
		).
		Complete(r)
}
