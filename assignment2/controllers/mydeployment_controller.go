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

	"github.com/go-logr/logr"
	velotiov1 "github.com/pankaj9310/k8s-dev-training/pankaj/assignment2/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// +kubebuilder:rbac:groups=velotio.pankaj.io,resources=mydeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=velotio.pankaj.io,resources=mydeployments/status,verbs=get;update;patch

func (r *MyDeploymentReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("mydeployment", req.NamespacedName)
	//your logic going here

	return ctrl.Result{}, nil
}

func (r *MyDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).
		For(&velotiov1.MyDeployment{}).
		Complete(r)
}
