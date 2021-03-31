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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
)

// CodeSanityReconciler reconciles a CodeSanity object
type CodeSanityReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type ResourceType string

const (
	JobResource    ResourceType = "Job"
	SanityResource ResourceType = "CodeSanity"
)

// +kubebuilder:rbac:groups=qa.test.com,resources=codesanities,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qa.test.com,resources=codesanities/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

func (r *CodeSanityReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("codesanity", req.NamespacedName)

	// Find out resource type from the request
	var sanity qav1.CodeSanity
	var job batchv1.Job

	resourceType := SanityResource

	if err := r.Get(ctx, req.NamespacedName, &sanity); err != nil {
		log.Info("Failed to get a CodeSanity. Trying to get a pod")

		if err = r.Get(ctx, req.NamespacedName, &job); err != nil {
			log.Error(err, "Failed to get a job")

			return ctrl.Result{}, err
		}

		resourceType = JobResource
	}

	switch resourceType {
	case SanityResource:
		return codesanity.CodeSanityRequestHandler(ctx, &sanity, r, r.Scheme, log)
	case JobResource:
		return jobs.JobRequestHandler(job, r, log)
	}

	return ctrl.Result{}, nil
}

func (r *CodeSanityReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&batchv1.Job{}, codesanity.JobOwnerKey, func(object runtime.Object) []string {
		job := object.(*batchv1.Job)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&qav1.CodeSanity{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
