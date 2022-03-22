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
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	apiv1 "github.com/hatred09/k8s-dev-training/assignment-4/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	jobOwnerKey   = ".metadata.controller"
	finalizerName = "my.domain.com/finalizer"
)

//+kubebuilder:rbac:groups=batch.my.domain.com,resources=scaler,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.my.domain.com,resources=scaler/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.my.domain.com,resources=scaler/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	fmt.Println("reconciler called")
	scaler := &apiv1.Scaler{}
	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		fmt.Println("could not fetch resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Finalizer Logic
	if scaler.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(scaler, finalizerName) {
			controllerutil.AddFinalizer(scaler, finalizerName)
			err := r.Update(ctx, scaler)
			if err != nil {
				return ctrl.Result{}, err
			}
			fmt.Println("finalizer added")
		}
	} else {
		if controllerutil.ContainsFinalizer(scaler, finalizerName) {
			err := r.RemoveDeployments(ctx, req)
			if err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(scaler, finalizerName)
			err = r.Update(ctx, scaler)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	replica := int32(scaler.Spec.Replicas)
	if scaler.Status.State == "" {
		scaler.Status.State = apiv1.Empty
	}

	switch scaler.Status.State {
	case apiv1.Empty:
		fmt.Println("create called")
		deploy := getDeploymentSpec(scaler.Spec.Name, scaler.Namespace, scaler.Spec.Image, scaler.Name, string(scaler.UID), scaler.APIVersion, &replica)
		err := r.Create(ctx, &deploy)
		if err != nil {
			return ctrl.Result{}, err
		}
		scaler.Status.State = apiv1.Created
		err = r.Status().Update(ctx, scaler)
		if err != nil {
			return ctrl.Result{}, err
		}
	default:
		deploy := v1.DeploymentList{}
		err := r.List(ctx, &deploy, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name})
		if err != nil {
			return ctrl.Result{}, err
		}
		if len(deploy.Items) == 0 {
			scaler.Status.State = apiv1.Empty
			err = r.Status().Update(ctx, scaler)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second * 5,
			}, nil
		}
		if scaler.Spec.Replicas != int(*deploy.Items[0].Spec.Replicas) {
			newDeploy := v1.Deployment{}
			getOpts := client.ObjectKey{
				Namespace: scaler.Namespace,
				Name:      scaler.Spec.Name,
			}
			err := r.Get(ctx, getOpts, &newDeploy)
			if err != nil {
				return ctrl.Result{}, err
			}
			newDeploy.Spec.Replicas = &replica
			r.Update(ctx, &newDeploy)
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	indexerFunc := func(rawObj client.Object) []string {
		deploy := rawObj.(*v1.Deployment)
		owner := metav1.GetControllerOf(deploy)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiv1.GroupVersion.String() || owner.Kind != "Scaler" {
			return nil
		}
		return []string{owner.Name}
	}

	err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1.Deployment{}, jobOwnerKey, indexerFunc)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.Scaler{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}

func (r *ScalerReconciler) RemoveDeployments(ctx context.Context, req ctrl.Request) error {
	var deploy v1.DeploymentList
	err := r.List(ctx, &deploy, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name})
	if err != nil {
		return err
	}
	if len(deploy.Items) == 0 {
		return nil
	}

	err = r.Delete(ctx, &deploy.Items[0])
	return err
}
