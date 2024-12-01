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

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

const indexedField = ".metadata.labels.job-name"
const finalizer = "mongodbs.db.velotio.com/dangling_cleaner"

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
			logger.Info("MongoDB resource not found, probably deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "get db instance")
		return ctrl.Result{}, err
	}
	// resource is marked for deletion, pending due to finalizer
	if !dbInstance.DeletionTimestamp.IsZero() {
		if err := r.cleanupNonLinkedResource(ctx, req.Namespace); err != nil {
			return ctrl.Result{}, err
		}
		if controllerutil.RemoveFinalizer(dbInstance, finalizer) {
			if err := r.Update(ctx, dbInstance); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	// create or update request, update the child deployment
	if err := r.upsertMongoPod(ctx, dbInstance); err != nil {
		logger.Error(err, "error while processing GVK")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if controllerutil.AddFinalizer(dbInstance, finalizer) {
		if err := r.Update(ctx, dbInstance); err != nil {
			return ctrl.Result{}, err
		}
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

func (r *MongoDBReconciler) cleanupNonLinkedResource(ctx context.Context, ns string) error {
	list := batchv1.JobList{}
	if err := r.List(ctx, &list, &client.MatchingFields{indexedField: "dummy-job"}, &client.ListOptions{Namespace: ns}); err != nil {
		return err
	}
	for _, v := range list.Items {
		if err := r.Delete(ctx, &v); err != nil {
			return err
		}
	}
	return nil
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
	p, ok := o.(*v1.Service)
	if !ok {
		return nil
	}
	indexValue, ok := p.Labels["job-name"]
	if !ok {
		return nil
	}
	return []string{indexValue}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &batchv1.Job{}, indexedField, indexer); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.MongoDB{}).
		WithEventFilter(eventProcessorPredicate()).
		Complete(r)
}
