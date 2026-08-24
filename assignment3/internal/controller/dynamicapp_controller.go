package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1alpha1 "github.com/Ajay-Raut-Rsystems/k8s-dev-training/assignment3/api/v1alpha1"
)

// DynamicAppReconciler reconciles a DynamicApp object
type DynamicAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.example.com,resources=dynamicapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=dynamicapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps;services,verbs=get;list;watch;create;update;patch;delete

func (r *DynamicAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. Fetch DynamicApp CR instance from local cache
	dynamicApp := &appsv1alpha1.DynamicApp{}
	err := r.Get(ctx, req.NamespacedName, dynamicApp)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("DynamicApp resource deleted. Skipping reconciliation.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var managedList []string

	// 2. Iterate through TargetResources specified in Spec
	for _, target := range dynamicApp.Spec.TargetResources {
		// Define dynamic object struct using unstructured type
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   target.Group,
			Version: target.Version,
			Kind:    target.Kind,
		})

		u.SetName(target.Name)
		u.SetNamespace(dynamicApp.Namespace)

		// Check if target resource exists in cluster
		err := r.Get(ctx, client.ObjectKeyFromObject(u), u)
		if errors.IsNotFound(err) {
			// Inject default content if creating a ConfigMap
			if target.Kind == "ConfigMap" {
				u.Object["data"] = map[string]interface{}{
					"managedBy": "Assignment3-DynamicAppController",
				}
			}

			// 3. Set OwnerReference for cascading deletion and event routing
			if err := controllerutil.SetControllerReference(dynamicApp, u, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}

			logger.Info("Creating Target Child Resource", "GVK", u.GroupVersionKind(), "Name", u.GetName())
			if err := r.Create(ctx, u); err != nil {
				return ctrl.Result{}, err
			}
		} else if err != nil {
			return ctrl.Result{}, err
		}

		managedList = append(managedList, fmt.Sprintf("%s/%s", target.Kind, target.Name))
	}

	// I face one error since I was trying to update an outdated version of the DynamicApp CR. To avoid this, I re-fetch the latest version of the DynamicApp before updating its status.
	// Error I have faced : Operation cannot be fulfilled on dynamicapps.apps.example.com "demo-dynamic-app": the object has been modified; please apply your changes to the latest version and try again
	// 4. Re-fetch the latest version of DynamicApp before updating status to avoid conflict errors
	latestApp := &appsv1alpha1.DynamicApp{}
	if err := r.Get(ctx, req.NamespacedName, latestApp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	latestApp.Status.Phase = "Running"
	latestApp.Status.ManagedResources = managedList
	if err := r.Status().Update(ctx, latestApp); err != nil {
		logger.Error(err, "Failed to update DynamicApp status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// Custom Predicate to filter unnecessary reconciliation events
func (r *DynamicAppReconciler) customPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Trigger reconcile ONLY if generation changes (filters out status updates)
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
	}
}

// SetupWithManager registers reconciler with Manager
func (r *DynamicAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.DynamicApp{}).
		WithEventFilter(r.customPredicate()). // Attach predicate filter
		Owns(&corev1.ConfigMap{}).            // Watch child ConfigMaps owned by this CR
		Complete(r)
}
