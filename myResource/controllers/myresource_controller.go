/*
Copyright 2021.

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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"myResource/helper"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	utilv2 "myResource/api/v2"
)

// MyResourceReconciler reconciles a MyResource object
type MyResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=util.my.domain,resources=myresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=util.my.domain,resources=myresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=util.my.domain,resources=myresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *MyResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// your logic here
	reqLogger := r.Log.WithValues("myResource", req.NamespacedName)
	instance := &utilv2.MyResource{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Println("resource not found in controller...")
			//	Object is found. Maybe deleted after reconcile request
			return ctrl.Result{}, nil
		}
		fmt.Println(err, ":: Unable to fetch procos instance...")
		//	object exists, but some other error
		return ctrl.Result{}, err
	}
	fmt.Println("instance status - ", instance.Status.JobState)
	if instance.Status.JobState == "" {
		instance.Status.JobState = utilv2.Uninitialised
		err := r.Status().Update(context.TODO(), instance)
		if err != nil {
			fmt.Println("Failed to update status of process")
			return ctrl.Result{}, err
		}
	}
	if instance.Status.JobState == "" {
		instance.Status.JobState = utilv2.Uninitialised
	}
	switch instance.Status.JobState {
	case utilv2.Uninitialised:
		fmt.Println("Initialising the job to execute after 5 minutes")
		instance.Status.JobState = utilv2.Pending
		instance.Spec.Schedule = time.Now().Add(5 * time.Minute).Format(time.RFC3339Nano)
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	case utilv2.Pending:
		fmt.Println("Pending State..")
		diff, err := helper.TimeLeft(instance.Spec.Schedule)
		if err != nil {
			reqLogger.Error(err, "Schedule parsing failure")
			return ctrl.Result{}, err
		}
		if diff > 0 {
			return ctrl.Result{RequeueAfter: diff * time.Second}, nil
		}
		fmt.Println("Executing the command")
		instance.Status.JobState = utilv2.Running
	case utilv2.Running:
		fmt.Println("Running State")
		pod := helper.SetPod(instance)
		err := ctrl.SetControllerReference(instance, pod, r.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}
		fmt.Println("Creating a POD")
		query := &corev1.Pod{}
		err = r.Get(context.TODO(), req.NamespacedName, query)
		if err != nil && errors.IsNotFound(err) {
			err = r.Create(context.TODO(), pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			reqLogger.Info("Pod Created successfully", "name", pod.Name)
			return ctrl.Result{}, nil
		} else if err != nil {
			// requeue with err
			reqLogger.Error(err, "cannot create pod")
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, nil
		}
	case utilv2.Finished:
		fmt.Println("Status Finished")
		return ctrl.Result{}, nil
	default:
		fmt.Println("Default state", instance.Status.JobState)
	}
	err = r.Status().Update(context.TODO(), instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		For(&utilv2.MyResource{}).
		Owns(&corev1.Pod{}).
		Complete(r)
	if err != nil {
		return err
	}
	return nil
}
