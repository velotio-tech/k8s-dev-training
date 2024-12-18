/*
Copyright 2024.

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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	creatorv1 "m3.io/creator/api/v1"
)

// ResourceCreatorReconciler reconciles a ResourceCreator object
type ResourceCreatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=creator.m3.io,resources=resourcecreators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=creator.m3.io,resources=resourcecreators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=creator.m3.io,resources=resourcecreators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceCreator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *ResourceCreatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	resourceCreator := &creatorv1.ResourceCreator{}
	err := r.Get(ctx, req.NamespacedName, resourceCreator)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch ResourceCreator")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, resourceSpec := range resourceCreator.Spec.Resources {
		logger.Info("reconciling resource", "resource", resourceSpec.Name)

		resource := &unstructured.Unstructured{}

		internalSpec := map[string]any{}
		err := json.Unmarshal(resourceSpec.Spec.Raw, &internalSpec)
		if err != nil {
			logger.Error(err, "unable to unmarshal resource spec", "resource", resourceSpec.Name)
			return ctrl.Result{}, err
		}

		// The Unstructured type represents a Kubernetes object that is not statically typed.
		// It uses a map[string]interface{} to store the object's data, including its GroupVersionKind (GVK), spec and metadata.
		// To avoid overwriting the GVK when setting the unstructured content,
		// we first set the spec as unstructured content and then explicitly set the GVK.
		// As all the methods that set fields in the Unstructured type basically inserts the key-value pair into the map[string]interface{}
		// and the GVK is not a separate field, it's just a part of the map[string]interface{}, we need to set the GVK explicitly after setting
		// the spec to prevent it from being overwritten.
		spec := map[string]any{"spec": internalSpec}
		resource.SetUnstructuredContent(spec)

		resource.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   resourceSpec.Group,
			Version: resourceSpec.Version,
			Kind:    resourceSpec.Kind,
		})

		resource.SetName(resourceSpec.Name)
		resource.SetNamespace(req.Namespace)

		ownerRefs := []metav1.OwnerReference{
			*metav1.NewControllerRef(resourceCreator, schema.GroupVersionKind{
				Group:   creatorv1.GroupVersion.Group,
				Version: creatorv1.GroupVersion.Version,
				Kind:    "ResourceCreator",
			}),
		}
		resource.SetOwnerReferences(ownerRefs)

		err = r.Get(ctx, client.ObjectKey{Name: resourceSpec.Name, Namespace: req.Namespace}, resource)
		if err != nil && client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch resource", "resource", resourceSpec.Name)
			return ctrl.Result{}, err
		}

		if err != nil && client.IgnoreNotFound(err) == nil {
			resource.SetResourceVersion("") // TODO: read why this is needed
			err = r.Create(ctx, resource)
			if err != nil {
				logger.Error(err, "unable to create resource", "resource", resourceSpec.Name)
				return ctrl.Result{}, err
			}
			logger.Info("created resource", "resource", resourceSpec.Name)
		} else {
			err = r.Update(ctx, resource)
			if err != nil {
				logger.Error(err, "unable to update resource", "resource", resourceSpec.Name)
				return ctrl.Result{}, err
			}
			logger.Info("updated resource", "resource", resourceSpec.Name)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceCreatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	builder := ctrl.NewControllerManagedBy(mgr).For(&creatorv1.ResourceCreator{})

	dynamicHandler := func(ctx context.Context, obj client.Object) []reconcile.Request {
		// Extract owner from the subresource and enqueue a reconcile request for the parent
		ownerRef := metav1.GetControllerOf(obj)
		if ownerRef != nil && ownerRef.Kind == "ResourceCreator" {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Name:      ownerRef.Name,
						Namespace: obj.GetNamespace(),
					},
				},
			}
		}
		return nil
	}

	gvkList := []schema.GroupVersionKind{
		{Group: "", Version: "v1", Kind: "Pod"},
		{Group: "", Version: "v1", Kind: "Service"},
		{Group: "apps", Version: "v1", Kind: "Deployment"},
		{Group: "apps", Version: "v1", Kind: "StatefulSet"},
		{Group: "apps", Version: "v1", Kind: "DaemonSet"},
		{Group: "batch", Version: "v1", Kind: "Job"},
	}

	for _, gvk := range gvkList {
		obj := &unstructured.Unstructured{}
		obj.SetGroupVersionKind(gvk)
		builder.Watches(obj, handler.EnqueueRequestsFromMapFunc(dynamicHandler))
	}

	return builder.Complete(r)
}
