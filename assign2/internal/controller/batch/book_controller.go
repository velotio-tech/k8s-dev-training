/*
Copyright 2023.

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
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	booksv1 "github.com/arijitnnayak/k8s-dev-training/api/batch/v1"
	v1 "github.com/arijitnnayak/k8s-dev-training/api/batch/v1"
)

// BookReconciler reconciles a Book object
type BookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	owner = "user-1"
)

//+kubebuilder:rbac:groups=books.assgn2.com,resources=books,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=books.assgn2.com,resources=books/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=books.assgn2.com,resources=books/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Book object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *BookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	// Fetch the Book resource

	desiredBookState := &booksv1.Book{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-book",
			Namespace: "default",
		},
		Spec: booksv1.BookSpec{
			Title:       "Desired Title",
			Author:      "Desired Author",
			Pages:       100,
			PublishDate: metav1.Now(),
			Reference: v1.ReferenceGVK{
				Group:   "example.com",
				Version: "v1",
				Kind:    "ExampleKind",
				Name:    "ExampleName",
			},
		},
		Status: booksv1.BookStatus{
			Available: true,
		},
	}

	currentBookState := &booksv1.Book{}
	if err := r.Get(ctx, req.NamespacedName, currentBookState); err != nil {
		if errors.IsNotFound(err) {
			// Book resource not found
			log.Info("book resource not found :  let's create a new one")
			owner := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example-deployment",
					Namespace: "default",
				},
			}
			// Set owner reference for the desiredBook
			if err := ctrl.SetControllerReference(owner, desiredBookState, r.Scheme); err != nil {
				log.Error(err, "Failed to set owner reference for Book resource")
				return reconcile.Result{}, err
			}

			if err := r.Create(ctx, desiredBookState); err != nil {
				log.Error(err, "failed to create Book resource")
				return ctrl.Result{}, err
			}
			log.Info("book created successfully")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Check if the book needs to be updated
	if !reflect.DeepEqual(currentBookState.Spec.Title, desiredBookState.Spec.Title) {
		// Perform update operation
		currentBookState.Spec.Pages = desiredBookState.Spec.Pages + 100%2
		if err := r.Status().Update(ctx, currentBookState, &client.SubResourceUpdateOptions{}); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("book updated successfully")
		return ctrl.Result{}, nil
	}

	// Check if the resource exists in the cluster but is not present in the desired state
	if currentBookState.ObjectMeta.DeletionTimestamp == nil {
		log.Info("current state indicates deletion. Deleting Book resource...")
		if err := r.Delete(ctx, currentBookState); err != nil {
			log.Error(err, "Failed to delete Book resource")
			return reconcile.Result{}, err
		}
		log.Info("book resource deleted successfully")
	}

	// Update the Available status based on the comparison
	available := reflect.DeepEqual(desiredBookState.Spec, desiredBookState.Spec)
	if currentBookState.Status.Available != available {
		log.Info("updating Available status in the Book resource status...")
		currentBookState.Status.Available = available
		if err := r.Status().Update(ctx, currentBookState); err != nil {
			log.Error(err, "failed to update Available status")
			return reconcile.Result{}, err
		}
		log.Info("available status updated successfully")
	}

	// Book is already in the desired state
	log.Info("Book is already in the desired state")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&booksv1.Book{}).
		Complete(r)
}
