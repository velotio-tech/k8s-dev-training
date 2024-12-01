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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "my.domain/ProcOS/api/v1"
	p_helper "my.domain/ProcOS/helpers/procos"
)

// ProcOSReconciler reconciles a ProcOS object
type ProcOSReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

const maxProcessAliveSeconds = 5
const maxProcessDeadSeconds = 5
const expectedProcOSDeploymentSize = 4
const expectedReplicaSize int32 = 4

var (
	// procosOwnerKey = ".metadata.controller"
	procosOwnerKey = ".metadata.namespace"
	apiGVStr       = batchv1.GroupVersion.String()
	process        *batchv1.ProcOS
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
	// log := r.Log.WithValues("process", req.NamespacedName)

	procos := &batchv1.ProcOS{}
	err := r.Get(context.TODO(), req.NamespacedName, procos)
	if err != nil {
		golog.Println(err, ":: Unable to fetch procos...")
	}

	procosDeps := p_helper.GetProcOSOwnedDeployments(ctx, r.Client, procos)
	if len(procosDeps.Items) >= expectedProcOSDeploymentSize {
		golog.Println("Cleaning up resources...")
		if err = r.cleanUpOwnedResources(context.Background(), procosDeps); err != nil {
			golog.Println(err, ":: Unable to clean-up owned resources...")
			return ctrl.Result{}, err
		}
		golog.Println("procos owned deployments are cleaned...")
	} else if len(procosDeps.Items) < expectedProcOSDeploymentSize {
		golog.Println("Creating procos owned deployments...")
		if err = createProcOSOwnedDeployments(ctx, r, procos); err != nil {
			golog.Println(err, ":: Unable to create expected procos deployments...")
			return ctrl.Result{}, err
		}
	}

	if err = r.updateProcOSDeploymentReplicas(ctx, procos); err != nil {
		golog.Println(err, ":: Unable to update replicas of deployments...")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(6 * time.Second)}, nil
}

func createProcOSOwnedDeployments(ctx context.Context, r *ProcOSReconciler, procos *batchv1.ProcOS) error {
	for len(p_helper.GetProcOSOwnedDeployments(ctx, r.Client, procos).Items) < expectedProcOSDeploymentSize {
		deployment := p_helper.GetProcOSDeploymentSpec("", procos.Namespace, procos)
		if err := r.Client.Create(ctx, deployment); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProcOSReconciler) cleanUpOwnedResources(ctx context.Context, procosDeps *appsv1.DeploymentList) error {
	deletePolicy := metav1.DeletePropagationBackground
	for _, dep := range procosDeps.Items {
		if err := r.Client.Delete(ctx, &dep, &client.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
			return err
		}
	}

	//	Delete pods with labels procos-dep
	labelSelector, _ := labels.ValidatedSelectorFromSet(map[string]string{"app": "procos-dep"})
	podList := &corev1.PodList{}
	r.Client.List(ctx, podList, &client.ListOptions{LabelSelector: labelSelector})
	fmt.Println("==========>pod label list: ", len(podList.Items))
	for _, pod := range podList.Items {
		// fmt.Println("====>pod name - ", pod.Name)
		if err := r.Client.Delete(ctx, &pod, &client.DeleteOptions{}); err != nil {
			golog.Println(err, ":: Unable to delete pod - ", pod.Name)
		}
		// pod.Status.Phase
		fmt.Println("pod - ", pod.Name, " deleted...")
	}

	return nil
}

//	Scales the replicaset of deployments.
func (r *ProcOSReconciler) updateProcOSDeploymentReplicas(ctx context.Context, procos *batchv1.ProcOS) error {
	procosDeps := p_helper.GetProcOSOwnedDeployments(ctx, r.Client, procos)
	for _, dep := range procosDeps.Items {
		if *dep.Spec.Replicas < expectedReplicaSize {
			*dep.Spec.Replicas *= 2 //	Make it double
			if err := r.Client.Update(ctx, &dep); err != nil {
				return err
			}
			golog.Printf("Scaled deployment %q to %d replicas.\n", dep.Name, *dep.Spec.Replicas)
		} else if *dep.Spec.Replicas > expectedReplicaSize {
			*dep.Spec.Replicas /= 2 //	Make it half
			if err := r.Client.Update(ctx, &dep); err != nil {
				return err
			}
			golog.Printf("Scaled deployment %q to %d replicas.\n", dep.Name, *dep.Spec.Replicas)
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProcOSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, procosOwnerKey, func(rawObj client.Object) []string {
		deployment := rawObj.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deployment)
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
