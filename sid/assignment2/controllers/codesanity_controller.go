/*


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
	"github.com/farkaskid/k8s-dev-training/assignment2/helpers/codesanity"
	"github.com/farkaskid/k8s-dev-training/assignment2/helpers/jobs"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
)

// CodeSanityReconciler reconciles a CodeSanity object
type CodeSanityReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=qa.test.com,resources=codesanities,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qa.test.com,resources=codesanities/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get
// +kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=cronjobs/status,verbs=get
func (r *CodeSanityReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("codesanity", req.NamespacedName)

	// Find out resource type from the request
	var sanity qav1.CodeSanity

	if err := r.Get(ctx, req.NamespacedName, &sanity); err != nil {
		log.Info("Failed to get a CodeSanity. Trying to get a pod")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return codesanity.CodeSanityRequestHandler(ctx, &sanity, r, r.Scheme, log)
}

func (r *CodeSanityReconciler) SetupWithManager(mgr ctrl.Manager) error {
	attachOwnerKey := func(object runtime.Object) []string {
		obj := object.(metav1.Object)
		owner := metav1.GetControllerOf(obj)
		if owner == nil {
			return nil
		}

		return []string{owner.Name}
	}

	if err := mgr.GetFieldIndexer().IndexField(&batchv1.Job{}, codesanity.JobOwnerKey, attachOwnerKey); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(&v1beta1.CronJob{}, codesanity.JobOwnerKey, attachOwnerKey); err != nil {
		return err
	}

	createFilter := func(e event.CreateEvent) bool {
		if _, ok := e.Object.(*qav1.CodeSanity); ok {
			return true
		}
		return false
	}

	deleteFilter := func(e event.DeleteEvent) bool {
		if _, ok := e.Object.(*qav1.CodeSanity); ok {
			return true
		}
		return false
	}

	updateFilter := func(e event.UpdateEvent) bool {
		if newJob, ok := e.ObjectNew.(*batchv1.Job); ok {
			return jobs.IsJobFinished(newJob)
		}

		if sanity, ok := e.ObjectNew.(*qav1.CodeSanity); ok {
			oldSanity, _ := e.ObjectOld.(*qav1.CodeSanity)

			if reflect.DeepEqual(sanity.Spec, oldSanity.Spec) {
				return false
			}
		}

		return true
	}

	p := predicate.Funcs{
		CreateFunc: createFilter,
		DeleteFunc: deleteFilter,
		UpdateFunc: updateFilter,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&qav1.CodeSanity{}).
		Owns(&batchv1.Job{}).
		Owns(&v1beta1.CronJob{}).
		WithEventFilter(p).
		Complete(r)
}
