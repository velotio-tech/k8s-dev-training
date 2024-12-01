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

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dbv1 "velotio.com/database/api/v1"
)

// MongoDBReconciler reconciles a MongoDB object
type MongoDBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=db.velotio.com,resources=mongodbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.velotio.com,resources=mongodbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.velotio.com,resources=mongodbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MongoDB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *MongoDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("reconciling...")

	dbInstance := &dbv1.MongoDB{}
	if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Name}, dbInstance); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("MongoDB resource not found")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "get db instance")
		return ctrl.Result{}, err
	}

	if err := r.upsertMongoPod(ctx, dbInstance); err != nil {
		logger.Error(err, "error while processing GVK")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	dbInstance.Status.Condition = "healthy"
	dbInstance.Status.Phase = "running"
	if err := r.Status().Update(ctx, dbInstance); err != nil {
		logger.Error(err, "update status")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("finished reconcilation...")
	return ctrl.Result{}, nil
}

func eventProcessorPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

func indexer(o client.Object) []string {
	p, ok := o.(*v1.Pod)
	if !ok {
		return nil
	}
	owner := metav1.GetControllerOf(p)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != dbv1.GroupVersion.Identifier() || owner.Kind != "MongoDB" {
		return nil
	}
	return []string{owner.Name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &v1.Pod{}, ".metadata.ownerReferences", indexer); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.MongoDB{}).
		WithEventFilter(eventProcessorPredicate()).
		Complete(r)
}
