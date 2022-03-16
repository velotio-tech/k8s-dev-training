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
	//"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
	apiGVStr    = v1.GroupVersion.String()
	secDel      *v1.SecretCreator
	isDelete    = false
	isUpdated   = false
	isCreated   = false
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
	var secretCreator v1.SecretCreator
	var secret, secrettocheck corev1.Secret
	if isDelete == true {
		if len(secDel.Spec.IncludeNamespaces) != 0 {
			for x := range secDel.Spec.IncludeNamespaces {
				if err := r.Get(ctx, types.NamespacedName{Name: secDel.Spec.SecretName, Namespace: secDel.Spec.IncludeNamespaces[x]}, &secrettocheck); err == nil {
					if err = r.DeleteSecret(ctx, secDel.Spec.IncludeNamespaces[x], secDel.Spec.SecretName, secrettocheck.Data); err != nil {
						return ctrl.Result{}, err
					}
				}
			}
		} else {
			namespaces, err := r.GetNamespaces(ctx)
			if err != nil {
				return ctrl.Result{}, err
			}
			for x := range namespaces.Items {
				if err = r.Get(ctx, types.NamespacedName{Name: secDel.Spec.SecretName,
					Namespace: namespaces.Items[x].ObjectMeta.Name}, &secrettocheck); err == nil {
					if err = r.DeleteSecret(ctx, namespaces.Items[x].ObjectMeta.Name, secDel.Spec.SecretName, secrettocheck.Data); err != nil {
						return ctrl.Result{}, err
					}
				}
			}
		}
		isDelete = false
		return ctrl.Result{}, nil
	}
	if err := r.Get(ctx, req.NamespacedName, &secretCreator); err != nil {
		r.Log.Error(err, "unable to fetch SecretCreator")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err := r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: secretCreator.Spec.SecretNamespace}, &secret); err != nil {
		r.Log.Error(err, "sceret not found", secretCreator.Spec.SecretName)
		return ctrl.Result{}, err
	}
	if len(secretCreator.Spec.IncludeNamespaces) == 0 && len(secretCreator.Spec.ExcludeNamespaces) == 0 {
		namespaces, err := r.GetNamespaces(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		for x := range namespaces.Items {
			if err := r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secrettocheck); err != nil {
				if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
					secretCreator, secret.Data); err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	} else if len(secretCreator.Spec.IncludeNamespaces) == 0 && len(secretCreator.Spec.ExcludeNamespaces) != 0 {
		namespaces, err := r.GetNamespaces(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		for x := range namespaces.Items {
			if !contains(secretCreator.Spec.ExcludeNamespaces, namespaces.Items[x].ObjectMeta.Name) {
				if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secrettocheck); err != nil {
					if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
						secretCreator, secret.Data); err != nil {
						return ctrl.Result{}, err
					}
				}
			} else {
				if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secrettocheck); err == nil {
					if err = r.DeleteSecret(ctx, namespaces.Items[x].ObjectMeta.Name, secretCreator.Spec.SecretName, secret.Data); err != nil {
						return ctrl.Result{}, err
					}
				}
			}
		}
	} else {
		namespaces, err := r.GetNamespaces(ctx)
		if err != nil {
			return ctrl.Result{}, err
		}
		for x := range namespaces.Items {
			if namespaces.Items[x].ObjectMeta.Name != secret.ObjectMeta.Namespace {
				if contains(secretCreator.Spec.IncludeNamespaces, namespaces.Items[x].ObjectMeta.Name) {
					if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secrettocheck); err != nil {
						if err = r.CreateSecret(ctx, namespaces.Items[x].ObjectMeta.Name,
							secretCreator, secret.Data); err != nil {
							return ctrl.Result{}, err
						}
					}
				} else {
					if err = r.Get(ctx, types.NamespacedName{Name: secretCreator.Spec.SecretName, Namespace: namespaces.Items[x].ObjectMeta.Name}, &secrettocheck); err == nil {
						if err = r.DeleteSecret(ctx, namespaces.Items[x].ObjectMeta.Name, secretCreator.Spec.SecretName, secret.Data); err != nil {
							return ctrl.Result{}, err
						}
					}
				}
			}
		}
	}
	time_now := &metav1.Time{Time: time.Now()}
	if secretCreator.Status.LastUpdatedTime == nil || isUpdated == true {
		secretCreator.Status.LastUpdatedTime = time_now
	}
	if secretCreator.Status.SecretNameSpace == "" || isUpdated == true {
		secretCreator.Status.SecretNameSpace = secretCreator.Spec.SecretNamespace
	}
	if isUpdated == true || isCreated == true {
		if err := r.Status().Update(ctx, &secretCreator); err != nil {
			r.Log.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
		isUpdated = false
		isCreated = false
	}
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
		if owner.APIVersion != apiGVStr || owner.Kind != "Secret" {
			return nil
		}
		fmt.Println(owner.Name)
		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	deleteFunction := func(e event.DeleteEvent) bool {
		if _, ok := e.Object.(*v1.SecretCreator); ok {
			secDel = e.Object.(*v1.SecretCreator)
			//r.Log.Info("Delete event for backup in backup controller")
			isDelete = true
			return true
		}
		return false
	}

	updateFunction := func(e event.UpdateEvent) bool {
		isUpdated = true
		return true
	}

	createFunction := func(e event.CreateEvent) bool {
		isCreated = true
		return true
	}

	p := predicate.Funcs{
		CreateFunc: createFunction,
		DeleteFunc: deleteFunction,
		UpdateFunc: updateFunction,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.SecretCreator{}).
		Owns(&corev1.Secret{}).
		WithEventFilter(p).
		Complete(r)
}
