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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	//"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/hatred09/k8s-dev-training/assignment-2/api/v1"
)

// SecretCreatorReconciler reconciles a SecretCreator object
type SecretCreatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var (
	jobOwnerKey = ".metadata.controller"
	sec         *v1.SecretCreator
)

//+kubebuilder:rbac:groups=secretcreator.example.com,resources=secretcreators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secretcreator.example.com,resources=secretcreators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secretcreator.example.com,resources=secretcreators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SecretCreator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SecretCreatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	fmt.Println("reconcile")
	secretCreator := v1.SecretCreator{}
	secretCreator.Status = v1.SecretCreatorStatus{}
	var secret, secretToCheck corev1.Secret

	namespaces, err := r.GetNamespaces(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	if sec != nil && sec.Status != secretCreator.Status && sec.Status.State == v1.Deleted {
		secretCreator.Status.State = v1.Deleted
	}

	switch secretCreator.Status.State {
	case v1.Deleted:
		fmt.Println("\n\n entry")
		if sec.Status.State == v1.Deleted {
			if len(sec.Spec.IncludeNamespaces) != 0 {
				for x := range sec.Spec.IncludeNamespaces {
					if err := r.Get(ctx, types.NamespacedName{Name: sec.Spec.SecretName, Namespace: sec.Spec.IncludeNamespaces[x]}, &secretToCheck); err == nil {
						if err = r.DeleteSecret(ctx, sec.Spec.IncludeNamespaces[x], sec.Spec.SecretName, secretToCheck.Data); err != nil {
							return ctrl.Result{}, err
						}
					}
				}
			} else {
				for x := range namespaces.Items {
					if err = r.Get(ctx, types.NamespacedName{Name: sec.Spec.SecretName,
						Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err == nil {
						if err = r.DeleteSecret(ctx, namespaces.Items[x].ObjectMeta.Name, sec.Spec.SecretName, secretToCheck.Data); err != nil {
							return ctrl.Result{}, err
						}
					}
				}
			}
			sec = nil
			secretCreator.Status.State = v1.Empty
		}
	default:
		if err := r.Get(ctx, req.NamespacedName, &secretCreator); err != nil {
			r.Log.Error(err, "unable to fetch SecretCreator")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		fmt.Println(req.NamespacedName, secretCreator, "==========================")

		fmt.Println("___________", namespaces.Items)
		if len(secretCreator.Spec.IncludeNamespaces) == 0 && len(secretCreator.Spec.ExcludeNamespaces) == 0 {
			fmt.Println("entry")
			for x := range namespaces.Items {
				if err := r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err != nil {
					if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
						secretCreator, secret.Data); err != nil {
						return ctrl.Result{}, err
					}
				}
			}
		} else if len(secretCreator.Spec.IncludeNamespaces) == 0 && len(secretCreator.Spec.ExcludeNamespaces) != 0 {
			fmt.Println("exclude")
			for x := range namespaces.Items {
				if !contains(secretCreator.Spec.ExcludeNamespaces, namespaces.Items[x].ObjectMeta.Name) {
					if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err != nil {
						if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
							secretCreator, secret.Data); err != nil {
							return ctrl.Result{}, err
						}
					}
				} else {
					if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err == nil {
						if err = r.DeleteSecret(ctx, namespaces.Items[x].ObjectMeta.Name, secretCreator.Spec.SecretName, secret.Data); err != nil {
							return ctrl.Result{}, err
						}
					}
				}
			}
		} else {
			for x := range namespaces.Items {
				fmt.Println("--------", x)
				if namespaces.Items[x].ObjectMeta.Name != secret.ObjectMeta.Namespace {
					if contains(secretCreator.Spec.IncludeNamespaces, namespaces.Items[x].ObjectMeta.Name) {
						if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err != nil {
							if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
								secretCreator, secret.Data); err != nil {
								return ctrl.Result{}, err
							}
						}
					} else {
						if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err == nil {
							if err = r.DeleteSecret(ctx, namespaces.Items[x].ObjectMeta.Name, secretCreator.Spec.SecretName, secret.Data); err != nil {
								return ctrl.Result{}, err
							}
						}
					}
				}
			}
		}
		if secretCreator.Status.State == "" {
			secretCreator.Status.State = v1.Created
		}
		if secretCreator.Status.LastUpdatedTime == nil {
			secretCreator.Status.LastUpdatedTime = &metav1.Time{Time: time.Now()}
		}
		if secretCreator.Status.SecretNameSpace == "" {
			secretCreator.Status.SecretNameSpace = secretCreator.Spec.SecretNamespace
		}
	}

	if err := r.Status().Update(ctx, &secretCreator); err != nil {
		return ctrl.Result{}, err
	}
	fmt.Println("++++++++++++++", ctrl.Result{}, secretCreator.Status, "+++++++++++++++")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretCreatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.Secret{}, jobOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner...
		job := rawObj.(*corev1.Secret)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}
		// ...make sure it's a CronJob...
		if owner.APIVersion != v1.GroupVersion.String() || owner.Kind != "Secret" {
			return nil
		}
		fmt.Println(owner.Name)
		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	deleteFunction := func(e event.DeleteEvent) bool {
		fmt.Println("\n\n delete event")
		if _, ok := e.Object.(*v1.SecretCreator); ok {
			sec = e.Object.(*v1.SecretCreator)
			sec.Status = v1.SecretCreatorStatus{}
			sec.Status.State = v1.Deleted
			return true
		}
		return false
	}

	p := predicate.Funcs{
		DeleteFunc: deleteFunction,
	}

	return ctrl.NewControllerManagedBy(mgr).For(&v1.SecretCreator{}).Owns(&corev1.Secret{}).WithEventFilter(p).Complete(r)
}
