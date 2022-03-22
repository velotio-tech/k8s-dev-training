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
	var secret, secretToCheck corev1.Secret
	if err := r.Get(ctx, req.NamespacedName, &secretCreator); err != nil {
		r.Log.Error(err, "unable to fetch SecretCreator")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err := r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: secretCreator.Spec.SecretNamespace}, &secret); err != nil {
		r.Log.Error(err, "secret not found", secretCreator.Spec.SecretName)
		return ctrl.Result{}, err
	}
	namespaces, err := r.GetNamespaces(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	if secretCreator.Status.State == "" {
		secretCreator.Status.State = v1.Empty
	}

	if sec.Status.State == v1.Delete {
		secretCreator.Status.State = v1.Delete
	}

	switch secretCreator.Status.State {
	case v1.Delete:
		if sec.Status.State == v1.Delete {
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
			sec.Status.State = v1.Empty
			secretCreator.Status.State = v1.Empty
			if err := r.Status().Update(ctx, &secretCreator); err != nil {
				r.Log.Error(err, "unable to update status")
				return ctrl.Result{}, err
			}
		}
	case v1.Empty, v1.Created:
		if len(secretCreator.Spec.IncludeNamespaces) == 0 && len(secretCreator.Spec.ExcludeNamespaces) == 0 {
			for x := range namespaces.Items {
				if err := r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secretToCheck); err != nil {
					if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
						secretCreator, secret.Data); err != nil {
						return ctrl.Result{}, err
					}
				}
			}
		} else if len(secretCreator.Spec.IncludeNamespaces) == 0 && len(secretCreator.Spec.ExcludeNamespaces) != 0 {
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
		secretCreator.Status.State = v1.Created
		if err := r.Status().Update(ctx, &secretCreator); err != nil {
			r.Log.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretCreatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	indexerFunc := func(rawObj client.Object) []string {
		secret := rawObj.(*corev1.Secret)
		owner := metav1.GetControllerOf(secret)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != v1.GroupVersion.String() || owner.Kind != "Secret" {
			return nil
		}
		return []string{owner.Name}
	}

	err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1.SecretCreator{}, jobOwnerKey, indexerFunc)
	if err != nil {
		return err
	}

	deleteFunction := func(e event.DeleteEvent) bool {
		if _, ok := e.Object.(*v1.SecretCreator); ok {
			sec = e.Object.(*v1.SecretCreator)
			sec.Status.State = v1.Delete
			return true
		}
		return false
	}

	updateFunction := func(e event.UpdateEvent) bool {
		return true
	}

	createFunction := func(e event.CreateEvent) bool {
		return true
	}

	p := predicate.Funcs{
		DeleteFunc: deleteFunction,
		CreateFunc: createFunction,
		UpdateFunc: updateFunction,
	}

	return ctrl.NewControllerManagedBy(mgr).For(&v1.SecretCreator{}).Owns(&corev1.Secret{}).WithEventFilter(p).Complete(r)
}
