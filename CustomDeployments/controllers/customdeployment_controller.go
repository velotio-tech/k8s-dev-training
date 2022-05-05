/*
Copyright 2022.

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
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	custdepv1 "paravkaushal.dev/CustomDeployments/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CustomDeploymentReconciler reconciles a CustomDeployment object
type CustomDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=custdep.paravkaushal.dev,resources=customdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=custdep.paravkaushal.dev,resources=customdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=custdep.paravkaushal.dev,resources=customdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CustomDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var log logr.Logger
	// TODO(user): your logic here
	customDeployment := custdepv1.CustomDeployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, &customDeployment); err != nil {
		log.Error(err, "cannot get customDeployment resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.cleanupOwnedResources(ctx, log, &customDeployment); err != nil {
		log.Error(err, "failed to clean up old Deployment resources for this customDeployment")
		return ctrl.Result{}, err
	}

	fmt.Println("Checking if an existing Deployment exists for this resource")
	deployment := apps.Deployment{}

	err := r.Client.Get(ctx, client.ObjectKey{Namespace: customDeployment.Namespace, Name: customDeployment.Spec.DeploymentName}, &deployment)

	if apierrors.IsNotFound(err) {
		fmt.Println("could not find existing deployment for CustomDeployment, creating one...")
		deployment := *buildDeployment(customDeployment)
		if err := r.Client.Create(ctx, &deployment); err != nil {
			log.Error(err, "failed to create Deployment resource")
			fmt.Println("created Deployment resource for customDeployment")
			return ctrl.Result{}, err
		}
	}
	if err != nil {
		log.Error(err, "failed to get Deployment for customDeployment resource")
		return ctrl.Result{}, err
	}

	expectedReplicas := int32(1)

	if customDeployment.Spec.Replicas != 0 {
		expectedReplicas = customDeployment.Spec.Replicas
	}
	if *deployment.Spec.Replicas != expectedReplicas {
		log.Info("updating replica count", "old_count", *deployment.Spec.Replicas, "new_count", expectedReplicas)

		deployment.Spec.Replicas = &expectedReplicas
		if err := r.Client.Update(ctx, &deployment); err != nil {
			log.Error(err, "failed to Deployment update replica count")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	log.Info("replica count up to date", "replica_count", *deployment.Spec.Replicas)

	log.Info("updating customDeployment resource status")
	customDeployment.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	if r.Client.Status().Update(ctx, &customDeployment); err != nil {
		log.Error(err, "failed to update customDeployment status")
		return ctrl.Result{}, err
	}

	log.Info("resource status synced")

	return ctrl.Result{}, nil
}

// cleanupOwnedResources will Delete any existing Deployment resources that
// were created for the given customDeployment that no longer match the
// customDeployment.spec.deploymentName field.
func (r *CustomDeploymentReconciler) cleanupOwnedResources(ctx context.Context, log logr.Logger, customDeployment *custdepv1.CustomDeployment) error {
	log.Info("finding existing Deployments for customDeployment resource")

	// List all deployment resources owned by this customDeployment
	var deployments apps.DeploymentList
	if err := r.List(ctx, &deployments, client.InNamespace(customDeployment.Namespace), client.MatchingFields(deploymentOwnerKey, customDeployment.Name)); err != nil {
		return err
	}

	deleted := 0
	for _, depl := range deployments.Items {
		if depl.Name == customDeployment.Spec.DeploymentName {
			// If this deployment's name matches the one on the customDeployment resource
			// then do not delete it.
			continue
		}

		if err := r.Client.Delete(ctx, &depl); err != nil {
			log.Error(err, "failed to delete Deployment resource")
			return err
		}
		deleted++
	}

	log.Info("finished cleaning up old Deployment resources", "number_deleted", deleted)

	return nil
}

func buildDeployment(customDeployment custdepv1.CustomDeployment) *apps.Deployment {
	deployment := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            customDeployment.Spec.DeploymentName,
			Namespace:       customDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&customDeployment, custdepv1.GroupVersion.WithKind("CustomDeployment"))},
		},
		Spec: apps.DeploymentSpec{
			Replicas: &customDeployment.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"custdep.paravkaushal.dev/deployment-name": customDeployment.Spec.DeploymentName,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"custdep.paravkaushal.dev/deployment-name": customDeployment.Spec.DeploymentName,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  "busybox",
							Image: "busybox:latest",
						},
					},
				},
			},
		},
	}
	return &deployment
}

var (
	deploymentOwnerKey = ".metadata.controller"
)

// SetupWithManager sets up the controller with the Manager.
func (r *CustomDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&custdepv1.CustomDeployment{}).
		Complete(r)
}
