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
	"fmt"

	webappv1 "demo/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AppdeployerReconciler reconciles a Appdeployer object
type AppdeployerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.velotio.ass2,resources=appdeployers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.velotio.ass2,resources=appdeployers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.velotio.ass2,resources=appdeployers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Appdeployer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile

var logger = ctrl.Log.WithName("controllerlog")

func (r *AppdeployerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	logger.Info("Reconciling Appdeployer custom resource..!")
	webapp := &webappv1.Appdeployer{}

	if err := r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: req.Name}, webapp); err != nil {
		// DELETE the service and deployment
		del_deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      webapp.Spec.Name + "-deployment",
				Namespace: webapp.Namespace,
			},
		}
		r.Delete(context.TODO(), del_deployment)
		fmt.Println("Deployment deleted successfully!!!")

		del_svc := &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      webapp.Spec.Name + "-service",
				Namespace: webapp.Namespace,
			},
		}
		r.Delete(context.TODO(), del_svc)
		fmt.Println("Service deleted successfully!!!")

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Println("Resource found..!", webapp.Name, webapp.Spec.Name, webapp.Spec.Replicas, webapp.Spec.ServiceType)
	// CREATE the service and deployment
	new_svc := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Spec.Name + "-service",
			Namespace: webapp.Namespace,
		},
		Spec: apiv1.ServiceSpec{
			Type: apiv1.ServiceType(webapp.Spec.ServiceType),
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(webapp.Spec.Port),
				},
			},
		},
	}
	//client.Client.Create(client, ctx, new_svc)
	if err := controllerutil.SetControllerReference(webapp, new_svc, r.Scheme); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	err := r.Create(context.TODO(), new_svc)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println("Service created successfully!!!")

	// Create a deployment
	var replicas int32 = int32(webapp.Spec.Replicas)
	new_deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      webapp.Spec.Name + "-deployment",
			Namespace: webapp.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"webapp": "appdeployer-lbl",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"webapp": "appdeployer-lbl",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Image: webapp.Spec.Image,
							Name:  webapp.Spec.Name + "-container",
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(webapp, new_deployment, r.Scheme); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	err = r.Create(context.TODO(), new_deployment)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Println("Deployment created successfully!!!")

	webapp.Status.AppProgress = "created"
	err = r.Status().Update(ctx, webapp)
	if err != nil {
		logger.Error(err, "Error while updating the status.")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppdeployerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Appdeployer{}).
		Complete(r)
}
