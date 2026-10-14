/*
Copyright 2026.

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

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	assignmentsv1 "k8s-dev-training.com/assignment-3/api/v1"
)

const (
	managedResourceIndex = "spec.managedResourceKeys"
	fieldManagerName     = "resourcemanager-controller"
)

// ResourceManagerReconciler reconciles a ResourceManager object
type ResourceManagerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=assignments.k8s-dev-training.com,resources=resourcemanagers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=assignments.k8s-dev-training.com,resources=resourcemanagers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=assignments.k8s-dev-training.com,resources=resourcemanagers/finalizers,verbs=update
// +kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.25.0/pkg/reconcile
func (r *ResourceManagerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	var rm assignmentsv1.ResourceManager

	if err := r.Get(ctx, req.NamespacedName, &rm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	previous := rm.Status.ManagedResources
	desiredRefs := make([]assignmentsv1.ManagedResourceReference, 0, len(rm.Spec.Resources))
	desiredKeys := make(map[string]struct{}, len(rm.Spec.Resources))

	for _, item := range rm.Spec.Resources {
		obj, ref, err := r.buildDesiredObject(ctx, &rm, item, req.Namespace)
		if err != nil {
			return ctrl.Result{}, err
		}

		key := resourceKey(ref)
		if _, exists := desiredKeys[key]; exists {
			return ctrl.Result{}, fmt.Errorf("duplicate managed resource: %s", key)
		}
		desiredKeys[key] = struct{}{}

		_, err = controllerutil.CreateOrUpdate(
			ctx,
			r.Client,
			obj,
			nil,
		)
		if err != nil {
			return ctrl.Result{}, err
		}
		desiredRefs = append(desiredRefs, ref)
	}

	for _, old := range previous {
		if _, stillDesired := desiredKeys[resourceKey(old)]; stillDesired {
			continue
		}

		if err := r.deleteManagedResource(ctx, &rm, old); err != nil {
			return ctrl.Result{}, err
		}
	}

	sort.Slice(desiredRefs, func(i, j int) bool {
		return resourceKey(desiredRefs[i]) < resourceKey(desiredRefs[j])
	})

	newStatus := rm.Status
	newStatus.ManagedResources = desiredRefs
	newStatus.Message = "All desired resources reconciled successfully"

	if !equality.Semantic.DeepEqual(rm.Status, newStatus) {
		rm.Status = newStatus
		if err := r.Status().Update(ctx, &rm); err != nil {
			return ctrl.Result{}, fmt.Errorf("update ResourceManager status: %w", err)
		}
	}
	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

func (r *ResourceManagerReconciler) buildDesiredObject(
	_ context.Context,
	rm *assignmentsv1.ResourceManager,
	item assignmentsv1.ManagedResource,
	namespace string,
) (*unstructured.Unstructured, assignmentsv1.ManagedResourceReference, error) {
	var empty assignmentsv1.ManagedResourceReference

	gv, err := schema.ParseGroupVersion(item.APIVersion)
	if err != nil {
		return nil, empty, fmt.Errorf("invalid apiVersion %q: %w", item.APIVersion, err)
	}
	if item.Kind == "" || item.Name == "" {
		return nil, empty, fmt.Errorf("kind and name are required")
	}

	gvk := gv.WithKind(item.Kind)

	obj := &unstructured.Unstructured{}
	if err := json.Unmarshal(item.Manifest.Raw, &obj.Object); err != nil {
		return nil, empty, fmt.Errorf(
			"decode manifest for %s/%s: %w", item.Kind, item.Name, err,
		)
	}

	obj.SetGroupVersionKind(gvk)
	obj.SetName(item.Name)
	obj.SetNamespace(namespace)

	// A namespaced owner can own a namespaced dependent in the same namespace.
	if err := controllerutil.SetControllerReference(rm, obj, r.Scheme); err != nil {
		return nil, empty, fmt.Errorf("set owner reference on %s: %w", item.Name, err)
	}

	ref := assignmentsv1.ManagedResourceReference{
		APIVersion: item.APIVersion,
		Kind:       item.Kind,
		Name:       item.Name,
	}

	return obj, ref, nil
}

func (r *ResourceManagerReconciler) deleteManagedResource(
	ctx context.Context,
	rm *assignmentsv1.ResourceManager,
	ref assignmentsv1.ManagedResourceReference,
) error {
	gv, err := schema.ParseGroupVersion(ref.APIVersion)
	if err != nil {
		return fmt.Errorf("invalid previous apiVersion %q: %w", ref.APIVersion, err)
	}

	gvk := gv.WithKind(ref.Kind)

	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	key := types.NamespacedName{
		Namespace: rm.GetNamespace(),
		Name:      ref.Name,
	}

	if err := r.Get(ctx, key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("get obsolete resource %s: %w", resourceKey(ref), err)
	}

	// Never delete an object that is no longer controlled by this CR.
	owner := metav1.GetControllerOf(obj)
	if owner == nil || owner.UID != rm.UID {
		return fmt.Errorf(
			"refusing to delete %s: it is not controlled by ResourceManager %s",
			resourceKey(ref), rm.Name,
		)
	}

	if err := r.Delete(ctx, obj); err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("delete obsolete resource %s: %w", resourceKey(ref), err)
	}

	return nil
}

func resourceKey(ref assignmentsv1.ManagedResourceReference) string {
	return ref.APIVersion + "|" + ref.Kind + "|" + ref.Name
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&assignmentsv1.ResourceManager{},
		managedResourceIndex,
		func(obj client.Object) []string {
			rm := obj.(*assignmentsv1.ResourceManager)
			keys := make([]string, 0, len(rm.Spec.Resources))

			for _, item := range rm.Spec.Resources {
				ref := assignmentsv1.ManagedResourceReference{
					APIVersion: item.APIVersion,
					Kind:       item.Kind,
					Name:       item.Name,
				}
				keys = append(keys, resourceKey(ref))
			}
			return keys
		},
	); err != nil {
		return err
	}

	specChangedPredicate := predicate.Funcs{
		CreateFunc:  func(event.CreateEvent) bool { return true },
		DeleteFunc:  func(event.DeleteEvent) bool { return true },
		GenericFunc: func(event.GenericEvent) bool { return true },
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldRM, okOld := e.ObjectOld.(*assignmentsv1.ResourceManager)
			newRM, okNew := e.ObjectNew.(*assignmentsv1.ResourceManager)
			if !okOld || !okNew {
				return true
			}

			return !equality.Semantic.DeepEqual(oldRM.Spec, newRM.Spec)
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&assignmentsv1.ResourceManager{}, builder.WithPredicates(specChangedPredicate)).
		Named("resourcemanager").
		Complete(r)
}
