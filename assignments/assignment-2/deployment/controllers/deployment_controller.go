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
	"math/rand"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "velotio.com/deployment/api/v1"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.velotio.com,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.velotio.com,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.velotio.com,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("[Reconcile][🎄]")

	deployment := &apiv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
		logger.Info("[🚨] Unable to fetch deployment")

		// unable to find deployment
		// check pod created by deployment
		// 	if pod is created or left from deployment delete it
		//  else do nothing
		var pods corev1.PodList
		if err := r.List(ctx, &pods, client.InNamespace(req.Namespace)); err != nil {
			logger.Error(err, "unable to fetch list of pods")
		}

		errcount := 0
		for i := 0; i < len(pods.Items); i++ {
			err := r.Delete(ctx, &pods.Items[i])
			logger.Info("[Deleting] deployment", "pods", pods.Items[i].ObjectMeta.Name)
			if err != nil {
				errcount++
			}
		}
		logger.Info("[Deleted] deployment pods")

		if errcount != 0 {
			return ctrl.Result{}, fmt.Errorf("pods deletetion failed")
		}

		return ctrl.Result{}, nil

	}

	// list all the pods
	var pods corev1.PodList
	if err := r.List(ctx, &pods, client.InNamespace(req.Namespace), client.MatchingLabels{"apps.velotio.com.deployment": deployment.Spec.Selector}); err != nil {
		logger.Error(err, "unable to fetch list of pods")
	}

	// only create remaing required pods
	provisionNewPods := deployment.Spec.Replicas - int32(len(pods.Items))

	// provision pods using corev1 api
	logger.Info("[Reconcile][🆕 provision pods]", "provision pods count", provisionNewPods)
	err := r.provisionPods(ctx, req, deployment, logger, provisionNewPods)
	if err != nil {
		return ctrl.Result{}, err
	}

	// update the status of ready replicas and replicas
	deployment.Status.ReadyReplicas = deployment.Spec.Replicas
	deployment.Status.Replicas = deployment.Spec.Replicas

	return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) provisionPods(ctx context.Context, req ctrl.Request, deployment *apiv1.Deployment, loggger logr.Logger, replicas int32) error {

	for i := 0; i < int(replicas); i++ {

		loggger.Info("Provision", "[🆗]", i)

		pod := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "apps.velotio.com.deployment.pod" + randomString(5),
				Namespace: req.Namespace,
				Labels: map[string]string{
					"apps.velotio.com.deployment": deployment.Spec.Selector,
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  deployment.Spec.Container.Name,
						Image: deployment.Spec.Container.Image,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(deployment.Spec.Container.Port),
							},
						},
					},
				},
			},
		}

		err := r.Create(ctx, pod)
		if err != nil {
			loggger.Error(err, "[🚨] error")
		}

	}

	return nil
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {

	rand.Seed(time.Now().UnixNano())

	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.Deployment{}).
		Complete(r)
}
