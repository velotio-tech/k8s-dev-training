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
	golog "log"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "my.domain/ProcOS/api/v1"
)

// ProcOSReconciler reconciles a ProcOS object
type ProcOSReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const maxProcessAliveSeconds = 5
const maxProcessDeadSeconds = 5

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = batchv1.GroupVersion.String()
	process     *batchv1.ProcOS
)

//+kubebuilder:rbac:groups=batch.my.domain,resources=procos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.my.domain,resources=procos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.my.domain,resources=procos/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProcOS object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ProcOSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	golog.Println("Reconciling procos...")

	instance := &batchv1.ProcOS{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			golog.Println("resource not found in controller...")
			//	Object is found. Maybe deleted after reconcile request
			return ctrl.Result{}, nil
		}
		golog.Println(err, ":: Unable to fetch procos instance...")
		//	object exists, but some other error
		return ctrl.Result{}, err
	}

	golog.Println("instance status - ", instance.Status.ProcessState)
	if instance.Status.ProcessState == "" {
		golog.Println("process status is empty string...")
		instance.Status.ProcessState = batchv1.Idle
		instance.Status.StartTime = &metav1.Time{Time: time.Now()}
		instance.Status.CompletionTime = &metav1.Time{Time: time.Now()}
		err := r.Status().Update(context.TODO(), instance)
		if err != nil {
			golog.Println("Failed to update status of process...")
			return ctrl.Result{}, err
		}
	}

	timeDiff := time.Now().Sub(instance.Status.StartTime.Time)
	completionTimeDiff := time.Now().Sub(instance.Status.CompletionTime.Time)
	fmt.Println("process alive time - ", timeDiff.Seconds())
	fmt.Println("process dead time - ", completionTimeDiff.Seconds())
	if instance.Status.ProcessState != batchv1.Zombie && instance.Status.ProcessState != batchv1.Stopped &&
		timeDiff.Seconds() > maxProcessAliveSeconds {
		golog.Println("Process - ", instance.GetName(), " is alive for a long time. Stopping the process' the process...")
		instance.Status.ProcessState = batchv1.Stopped
		instance.Status.CompletionTime = &metav1.Time{Time: time.Now()}

		return ctrl.Result{}, r.Status().Update(context.TODO(), instance)
	} else if (instance.Status.ProcessState == batchv1.Stopped ||
		instance.Status.ProcessState == batchv1.Zombie) &&
		timeDiff.Seconds() > maxProcessDeadSeconds {
		golog.Println("Process is dead for a long time. Recreating process...")
		instance.Status.ProcessState = batchv1.Idle
		instance.Status.StartTime = &metav1.Time{Time: time.Now()}

		return ctrl.Result{}, r.Status().Update(context.TODO(), instance)
	}

	switch instance.Status.ProcessState {

	case batchv1.Idle:
		fmt.Println("Scheduling Idle process - ", instance.GetName(), " to run.")
		instance.Status.ProcessState = batchv1.Running
		break

	case batchv1.InterruptableSleep:
		fmt.Println("Interrupting process - ", instance.GetName())
		instance.Status.ProcessState = batchv1.Running
		break

	case batchv1.Stopped:
		fmt.Println("Zombie'ing' the stopped process - ", instance.GetName())
		instance.Status.ProcessState = batchv1.Zombie
		break

	case batchv1.Zombie:
		fmt.Println("Process - ", instance.GetName(), " is zombie state. Recreating the process...")
		instance.Status.ProcessState = batchv1.Idle
		instance.Status.StartTime = &metav1.Time{Time: time.Now()}
		instance.Status.CompletionTime = &metav1.Time{Time: time.Now()}
		break

	default:
		golog.Println("process status - ", instance.Status.ProcessState)
		break
	}

	err = r.Status().Update(context.TODO(), instance)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(5 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProcOSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &batchv1.ProcOS{}, jobOwnerKey, func(rawObj client.Object) []string {
		process := rawObj.(*batchv1.ProcOS)
		owner := metav1.GetControllerOf(process)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind == "ProcOS" {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.ProcOS{}).
		WithEventFilter(getEventPredicates()).
		Complete(r)
}

func getEventPredicates() predicate.Predicate {
	pred := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			proc := e.Object.(*batchv1.ProcOS)
			fmt.Println("Process - ", proc.GetName(), " is created")
			return true
		},
		UpdateFunc: func(ue event.UpdateEvent) bool {
			proc := ue.ObjectOld.(*batchv1.ProcOS)
			fmt.Println("Process - ", proc.GetName(), " was updated.")
			return true
		},
		DeleteFunc: func(de event.DeleteEvent) bool {
			if _, ok := de.Object.(*batchv1.ProcOS); ok {
				process = de.Object.(*batchv1.ProcOS)
				golog.Println("deleted object is cached...")
			}
			return true
		},
	}

	return pred
}
