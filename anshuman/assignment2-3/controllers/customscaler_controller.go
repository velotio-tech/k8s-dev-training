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
	"errors"
	"fmt"
	logger "log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	batchv1 "domain.example.com/project/api/v1"
	"domain.example.com/project/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CustomScalerReconciler reconciles a CustomScaler object
type CustomScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	jobOwnerKey   = ".metadata.controller"
	finalizerName = "batch.domain.example.com/finalizer"
)

//+kubebuilder:rbac:groups=batch.domain.example.com,resources=customscalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.domain.example.com,resources=customscalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.domain.example.com,resources=customscalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomScaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *CustomScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger.Println("reconcile called")

	customScaler := &batchv1.CustomScaler{}
	err := r.Get(context.TODO(), req.NamespacedName, customScaler)
	if err != nil {
		logger.Println("Could not fetch resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if customScaler.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(customScaler, finalizerName) {
			controllerutil.AddFinalizer(customScaler, finalizerName)
			err := r.Update(ctx, customScaler)
			if err != nil {
				return ctrl.Result{}, err
			}
			logger.Println("finalizer added")
		}
	} else {
		if controllerutil.ContainsFinalizer(customScaler, finalizerName) {
			err := r.deleteResources(ctx, req, customScaler)
			if err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(customScaler, finalizerName)
			err = r.Update(ctx, customScaler)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if customScaler.Status.State == "" {
		customScaler.Status.State = batchv1.Uninitialized
	}

	switch customScaler.Status.State {
	case batchv1.Uninitialized:
		logger.Println("Creating pods for CustomScaler")
		for i := 0; i < customScaler.Spec.Replicas; i++ {
			podName := fmt.Sprintf("%s-%d", req.Name, i)
			pod := utils.GetPodSpec(podName, customScaler.Spec.Image, req.Namespace, customScaler.Name, string(customScaler.UID), customScaler.APIVersion)
			err := r.Create(context.TODO(), &pod)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		logger.Println("Pods created")
		customScaler.Status.State = batchv1.Created
		err := r.Status().Update(ctx, customScaler)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	// case batchv1.Error:
	// 	for _, pod := range customScaler.Status.Pods {
	// 		err := r.Delete(context.TODO(), &pod)
	// 		if err != nil {
	// 			logger.Println("error while deleting pod")
	// 			return ctrl.Result{}, err
	// 		}
	// 	}
	// 	logger.Println("Deleted pods due to error state, trying to recover")
	// 	customScaler.Status.State = batchv1.Uninitialized
	// 	err := r.Status().Update(context.TODO(), &customScaler)
	// 	if err != nil {
	// 		return ctrl.Result{}, err
	// 	}
	case batchv1.Created:
		requirement, _ := labels.NewRequirement("owner", selection.Equals, []string{customScaler.Name})
		sel := labels.NewSelector()
		sel.Add(*requirement)
		pods := v1.PodList{}
		err := r.List(ctx, &pods, &client.ListOptions{
			LabelSelector: sel,
			Namespace:     customScaler.Namespace,
		})
		if err != nil {
			return ctrl.Result{}, err
		}
		noOfPods := len(pods.Items)
		if noOfPods != customScaler.Spec.Replicas {
			if customScaler.Spec.Replicas < 1 {
				return ctrl.Result{}, errors.New("invalid config")
			}
			if noOfPods > customScaler.Spec.Replicas {
				difference := noOfPods - customScaler.Spec.Replicas
				for i := 0; i < difference; i++ {
					r.Delete(ctx, &pods.Items[noOfPods-i-1])
				}
			} else {
				difference := customScaler.Spec.Replicas - noOfPods
				for i := 0; i < difference; i++ {
					podName := fmt.Sprintf("%s-%d", req.Name, i+noOfPods)
					sp := utils.GetPodSpec(podName, customScaler.Spec.Image, req.Namespace, customScaler.Name, string(customScaler.UID), customScaler.APIVersion)
					r.Create(ctx, &sp, &client.CreateOptions{})
				}
			}
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	indexerFunc := func(rawObj client.Object) []string {
		customScaler := rawObj.(*batchv1.CustomScaler)
		owner := metav1.GetControllerOf(customScaler)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != batchv1.GroupVersion.String() || owner.Kind != "customScaler" {
			return nil
		}
		return []string{owner.Name}
	}

	err := mgr.GetFieldIndexer().IndexField(context.Background(), &batchv1.CustomScaler{}, jobOwnerKey, indexerFunc)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CustomScaler{}).
		Owns(&v1.Pod{}).
		Complete(r)
}

func (r *CustomScalerReconciler) deleteResources(ctx context.Context, req ctrl.Request, customScaler *batchv1.CustomScaler) error {
	var podList v1.PodList
	requirement, _ := labels.NewRequirement("owner", selection.Equals, []string{customScaler.Name})
	sel := labels.NewSelector()
	sel.Add(*requirement)
	err := r.List(ctx, &podList, &client.ListOptions{
		LabelSelector: sel,
		Namespace:     customScaler.Namespace,
	})
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		r.Delete(ctx, &pod, &client.DeleteOptions{})
	}
	return nil
}
