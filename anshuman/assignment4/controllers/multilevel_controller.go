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
	logger "log"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "my.example.com/assignment/api/v1"
)

// MultiLevelReconciler reconciles a MultiLevel object
type MultiLevelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	jobOwnerKey   = ".metadata.controller"
	finalizerName = "my.example.com/finalizer"
)

//+kubebuilder:rbac:groups=batch.my.example.com,resources=multilevels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.my.example.com,resources=multilevels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.my.example.com,resources=multilevels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MultiLevel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MultiLevelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger.Println("reconciler called")
	multiLevel := &batchv1.MultiLevel{}

	err := r.Get(ctx, req.NamespacedName, multiLevel)
	if err != nil {
		logger.Println("could not fetch resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Finalizer Logic
	if multiLevel.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(multiLevel, finalizerName) {
			controllerutil.AddFinalizer(multiLevel, finalizerName)
			err := r.Update(ctx, multiLevel)
			if err != nil {
				return ctrl.Result{}, err
			}
			logger.Println("finalizer added")
		}
	} else {
		if controllerutil.ContainsFinalizer(multiLevel, finalizerName) {
			err := r.cleanup(ctx, req)
			if err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(multiLevel, finalizerName)
			err = r.Update(ctx, multiLevel)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	repl := int32(multiLevel.Spec.Replicas)
	if multiLevel.Status.State == "" {
		logger.Println("create called")
		depl := getDeploymentSpec(multiLevel.Spec.Name, multiLevel.Namespace, multiLevel.Spec.Image, multiLevel.Name, string(multiLevel.UID), multiLevel.APIVersion, &repl)
		err := r.Create(ctx, &depl)
		if err != nil {
			return ctrl.Result{}, err
		}
		multiLevel.Status.State = "Created"
		err = r.Status().Update(ctx, multiLevel)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		depl := v1.DeploymentList{}
		err := r.List(ctx, &depl, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name})
		if err != nil {
			return ctrl.Result{}, err
		}
		if len(depl.Items) == 0 {
			multiLevel.Status.State = ""
			err = r.Status().Update(ctx, multiLevel)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second * 2,
			}, nil
		}
		if multiLevel.Spec.Replicas != int(*depl.Items[0].Spec.Replicas) {
			newDepl := v1.Deployment{}
			getOpts := client.ObjectKey{
				Namespace: multiLevel.Namespace,
				Name:      multiLevel.Spec.Name,
			}
			err := r.Get(ctx, getOpts, &newDepl)
			if err != nil {
				return ctrl.Result{}, err
			}
			newDepl.Spec.Replicas = &repl
			r.Update(ctx, &newDepl)
		}
	}

	return ctrl.Result{}, nil
}

func (r *MultiLevelReconciler) cleanup(ctx context.Context, req ctrl.Request) error {
	var depl v1.DeploymentList
	err := r.List(ctx, &depl, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name})
	if err != nil {
		return err
	}
	if len(depl.Items) == 0 {
		return nil
	}

	err = r.Delete(ctx, &depl.Items[0])
	return err
}

func getDeploymentSpec(name, namespace, image, ownerName, ownerUID, ownerAPIVersion string, replicas *int32) v1.Deployment {
	t := true
	depl := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: ownerAPIVersion,
				Kind:       "MultiLevel",
				Name:       ownerName,
				UID:        types.UID(ownerUID),
				Controller: &t,
			}},
			Namespace: namespace,
			Labels:    map[string]string{"app": "deployment"},
		},
		Spec: v1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:         name,
					GenerateName: name,
					Namespace:    name,
					Labels:       map[string]string{"app": "deployment"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  name,
						Image: image,
					}},
				},
			},
		},
	}
	return depl
}

// SetupWithManager sets up the controller with the Manager.
func (r *MultiLevelReconciler) SetupWithManager(mgr ctrl.Manager) error {

	indexerFunc := func(rawObj client.Object) []string {
		depl := rawObj.(*v1.Deployment)
		owner := metav1.GetControllerOf(depl)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != batchv1.GroupVersion.String() || owner.Kind != "MultiLevel" {
			return nil
		}
		return []string{owner.Name}
	}

	err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1.Deployment{}, jobOwnerKey, indexerFunc)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.MultiLevel{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}
