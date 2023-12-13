/*
Copyright 2023.

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

package pens

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pensv1 "github.com/arijitnnayak/k8s-dev-training/api/pens/v1"
)

// PenReconciler reconciles a Pen object
type PenReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pens.assgn2.com,resources=pens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pens.assgn2.com,resources=pens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pens.assgn2.com,resources=pens/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pen object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// Fetch the Pen resource
	pen := &pensv1.Pen{}
	if err := r.Get(ctx, req.NamespacedName, pen); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pen-deployment",
			Namespace: req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			// Define Deployment spec...
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "pen-app",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "pen-app",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "pen-container",
							Image: "ngnix",
						},
					},
				},
			},
		},
	}

	// Check if the Pen resource needs to be created
	if pen.ObjectMeta.DeletionTimestamp == nil && !pen.Status.Available {
		// Create logic for the Pen resource
		// Logic to create a Deployment

		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pen-job",
				Namespace: req.Namespace,
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "pen-job-container",
								Image: "ngnix",
							},
						},
						RestartPolicy: corev1.RestartPolicyOnFailure,
					},
				},
			},
		}
		if err := ctrl.SetControllerReference(pen, job, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, job); err != nil {
			return ctrl.Result{}, err
		}

		// Update Status indicating resource availability
		pen.Status.Available = true
		if err := r.Status().Update(ctx, pen); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Check if the Pen resource needs to be updated
	existingPen := &pensv1.Pen{}
	if err := r.Get(ctx, req.NamespacedName, existingPen); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Perform comparison logic between existing state and pen.Spec
	if existingPen.Spec.Name != pen.Spec.Name ||
		existingPen.Spec.Color != pen.Spec.Color {
		existingPen.Spec = pen.Spec
		if err := r.Update(ctx, existingPen); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Check if the Pen resource needs to be deleted
	if pen.ObjectMeta.DeletionTimestamp != nil {
		// Delete logic for the Pen resource
		if err := r.Delete(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}

		// Remove the finalizer and delete the resource
		ctrl.Log.Info("Deleting Pen resource...")
		if err := r.Delete(ctx, pen); err != nil {
			return ctrl.Result{}, err
		}

		// Exit the reconciliation loop
		return ctrl.Result{}, nil
	}

	// Fetch child resources (Pods) associated with the Pen using OwnerReference
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(req.Namespace), client.MatchingFields(map[string]string{
		"metadata.ownerReferences": pen.GetName(),
	})); err != nil {
		return ctrl.Result{}, err
	}

	// Perform operations on child resources (Pods) associated with the Pen
	for _, pod := range podList.Items {
		ctrl.Log.Info(fmt.Sprintf("Pod %s associated with Pen %s", pod.GetName(), pen.GetName()))
	}

	// Use field indexer to get resources (e.g., Jobs) without OwnerReference
	jobList := &batchv1.JobList{}
	if err := r.List(ctx, jobList, client.InNamespace(req.Namespace), client.MatchingFields{
		".metadata.ownerReferences": "",
	}); err != nil {
		return ctrl.Result{}, err
	}

	// Perform operations on resources (Jobs) without OwnerReference
	for _, job := range jobList.Items {
		ctrl.Log.Info(fmt.Sprintf("Job %s has no OwnerReference", job.GetName()))
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *PenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pensv1.Pen{}).
		Complete(r)
}
