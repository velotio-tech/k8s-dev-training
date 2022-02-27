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
	logger "log"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "domain.example.com/project/api/v1"
	"domain.example.com/project/utils"
)

// CustomScalerReconciler reconciles a CustomScaler object
type CustomScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

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
	_ = log.FromContext(ctx)

	var customScaler batchv1.CustomScaler
	err := r.Get(context.TODO(), req.NamespacedName, &customScaler)
	if err != nil {
		logger.Println("Could not fetch resource")
		return ctrl.Result{}, err
	}

	logger.Println("Resource name: ", customScaler.Name)

	if customScaler.Status.State == "" {
		customScaler.Status.State = batchv1.Uninitialized
	}

	switch customScaler.Status.State {
	case batchv1.Uninitialized:
		logger.Println("Creating pods for CustomScaler")
		for i := 0; i < customScaler.Spec.Replicas; i++ {
			podName := fmt.Sprintf("%s-%d", req.Name, i)
			pod := utils.GetPodSpec(podName, customScaler.Spec.Image)
			err := r.Create(context.TODO(), &pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			customScaler.Status.Pods = append(customScaler.Status.Pods, pod)
		}
		logger.Println("Pods created")
		customScaler.Status.State = batchv1.Created
		err := r.Status().Update(context.TODO(), &customScaler)
		if err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	case batchv1.Error:
		for _, pod := range customScaler.Status.Pods {
			err := r.Delete(context.TODO(), &pod)
			if err != nil {
				logger.Println("error while deleting pod")
				return ctrl.Result{}, err
			}
		}
		logger.Println("Deleted pods due to error state, trying to recover")
		customScaler.Status.State = batchv1.Uninitialized
		err := r.Status().Update(context.TODO(), &customScaler)
		if err != nil {
			return ctrl.Result{}, err
		}
	case batchv1.Created:
		for _, pod := range customScaler.Status.Pods {
			if pod.Status.Phase != "Running" {
				customScaler.Status.State = batchv1.Error
				err := r.Status().Update(context.TODO(), &customScaler)
				if err != nil {
					return ctrl.Result{}, err
				}
				break
			}
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CustomScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CustomScaler{}).
		Complete(r)
}
