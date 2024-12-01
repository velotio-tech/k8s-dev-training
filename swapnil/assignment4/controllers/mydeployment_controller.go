// /*

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package controllers

import (
	"context"
	"log"

	"github.com/go-logr/logr"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	mydeplyv1 "swapnil/k8s-dev-training/swapnil/assignment4/api/v1"
)

// MyDeploymentReconciler reconciles a MyDeployment object
type MyDeploymentReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var (
	deploymentOwnerKey = ".metadata.controller"
)

// kubebuilder:rbac:groups=velotio.pankaj.io,resources=mydeployments,verbs=get;list;watch;create;update;patch;delete
// kubebuilder:rbac:groups=velotio.pankaj.io,resources=mydeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *MyDeploymentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("mydeployment", req.NamespacedName)

	log.Info("Fetching MyDeployment resources.")
	myDeployment := mydeplyv1.MyDeployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, &myDeployment); err != nil {
		log.Error(err, "failed to get MyDeployment resource")
		// Ignore NotFound errors as they will be retried automatically if the
		// resource is created in future.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if err := r.cleanupOwnedResources(ctx, log, &myDeployment); err != nil {
		log.Error(err, "failed to clean up old Deployment resources for this MyDeployment")
		return ctrl.Result{}, err
	}

	log = log.WithValues("deployment_name", myDeployment.Spec.DeploymentName)

	log.Info("checking if an existing Deployment exists for this resource")
	deployment := apps.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: myDeployment.Namespace, Name: myDeployment.Spec.DeploymentName}, &deployment)
	if apierrors.IsNotFound(err) {
		log.Info("could not find existing Deployment for MyDeployment, creating one...")

		deployment = *buildDeployment(myDeployment)
		if err := r.Client.Create(ctx, &deployment); err != nil {
			log.Error(err, "failed to create Deployment resource")
			return ctrl.Result{}, err
		}

		r.Recorder.Eventf(&myDeployment, core.EventTypeNormal, "Created", "Created deployment %q", deployment.Name)
		log.Info("created Deployment resource for MyDeployment")
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "failed to get Deployment for MyDeployment resource")
		return ctrl.Result{}, err
	}

	log.Info("existing Deployment resource already exists for MyDeployment, checking replica count")

	expectedReplicas := int32(1)
	if myDeployment.Spec.Replicas != nil {
		expectedReplicas = *myDeployment.Spec.Replicas
	}
	if *deployment.Spec.Replicas != expectedReplicas {
		log.Info("updating replica count", "old_count", *deployment.Spec.Replicas, "new_count", expectedReplicas)

		deployment.Spec.Replicas = &expectedReplicas
		if err := r.Client.Update(ctx, &deployment); err != nil {
			log.Error(err, "failed to Deployment update replica count")
			return ctrl.Result{}, err
		}

		r.Recorder.Eventf(&myDeployment, core.EventTypeNormal, "Scaled", "Scaled deployment %q to %d replicas", deployment.Name, expectedReplicas)

		return ctrl.Result{}, nil
	}

	log.Info("replica count up to date", "replica_count", *deployment.Spec.Replicas)

	log.Info("updating MyDeployment resource status")
	myDeployment.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	if r.Client.Status().Update(ctx, &myDeployment); err != nil {
		log.Error(err, "failed to update MyDeployment status")
		return ctrl.Result{}, err
	}

	log.Info("resource status synced")

	return ctrl.Result{}, nil
}

// cleanupOwnedResources will Delete any existing Deployment resources that
// were created for the given MyDeployment that no longer match the
// myDeployment.spec.deploymentName field.
func (r *MyDeploymentReconciler) cleanupOwnedResources(ctx context.Context, log logr.Logger, myDeployment *mydeplyv1.MyDeployment) error {
	log.Info("finding existing Deployments for MyDeployment resource")

	// List all deployment resources owned by this MyDeployment
	var deployments apps.DeploymentList
	if err := r.List(ctx, &deployments, client.InNamespace(myDeployment.Namespace), client.MatchingField(deploymentOwnerKey, myDeployment.Name)); err != nil {
		return err
	}

	deleted := 0
	for _, depl := range deployments.Items {
		if depl.Name == myDeployment.Spec.DeploymentName {
			// If this deployment's name matches the one on the MyDeployment resource
			// then do not delete it.
			continue
		}

		if err := r.Client.Delete(ctx, &depl); err != nil {
			log.Error(err, "failed to delete Deployment resource")
			return err
		}

		r.Recorder.Eventf(myDeployment, core.EventTypeNormal, "Deleted", "Deleted deployment %q", depl.Name)
		deleted++
	}

	log.Info("finished cleaning up old Deployment resources", "number_deleted", deleted)

	return nil
}

func buildDeployment(myDeployment mydeplyv1.MyDeployment) *apps.Deployment {
	deployment := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            myDeployment.Spec.DeploymentName,
			Namespace:       myDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&myDeployment, mydeplyv1.GroupVersion.WithKind("MyDeployment"))},
		},
		Spec: apps.DeploymentSpec{
			Replicas: myDeployment.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"example-controller.jetstack.io/deployment-name": myDeployment.Spec.DeploymentName,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"example-controller.jetstack.io/deployment-name": myDeployment.Spec.DeploymentName,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  "pankaj-myjob",
							Image: "docker.io/pankaj9310/myjob:latest",
							Command: []string{
								"./app",
							},
						},
					},
				},
			},
		},
	}
	return &deployment
}

func (r *MyDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&apps.Deployment{}, deploymentOwnerKey, func(rawObj runtime.Object) []string {
		// grab the Deployment object, extract the owner...
		depl := rawObj.(*apps.Deployment)
		owner := metav1.GetControllerOf(depl)
		if owner == nil {
			return nil
		}
		// ...make sure it's a MyDeployment...
		if owner.APIVersion != mydeplyv1.GroupVersion.String() || owner.Kind != "MyDeployment" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	checkNamespace := func(e event.CreateEvent) bool {
		if e.Meta.GetNamespace() != "kube-system" {
			log.Printf("error: allowed namespace is kube-system, found %s", e.Meta.GetNamespace())
			return false
		}
		return true
	}

	p := predicate.Funcs{
		CreateFunc: checkNamespace,
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&mydeplyv1.MyDeployment{}).
		Owns(&apps.Deployment{}).
		WithEventFilter(p).
		Complete(r)
}
