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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8sassignment2v1 "my.domain/k8s-assignment-2/api/v1"
)

// FlightTicketReconciler reconciles a FlightTicket object
type FlightTicketReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8s-assignment-2.my.domain,resources=flighttickets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s-assignment-2.my.domain,resources=flighttickets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s-assignment-2.my.domain,resources=flighttickets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FlightTicket object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *FlightTicketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling FlightTicket custom resource")

	ft := &k8sassignment2v1.FlightTicket{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, ft)

	ft.Status.BookingStatus = "DONE"
	ft.Status.Fare = 5000
	if err := r.Status().Update(ctx, ft); err != nil {
		log.Error(err, "unable to update flightticket status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FlightTicketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sassignment2v1.FlightTicket{}).
		Complete(r)
}
