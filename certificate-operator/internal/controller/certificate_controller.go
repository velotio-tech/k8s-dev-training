/*
Copyright 2024.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	certv1 "github.com/jshiwamv/k8s-dev-training/certificate-operator/api/v1"
)

// CertificateReconciler reconciles a Certificate object
type CertificateReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cert.example.com,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert.example.com,resources=certificates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cert.example.com,resources=certificates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Certificate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *CertificateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	cert := &certv1.Certificate{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, cert)
	if err != nil {
		return ctrl.Result{}, err
	}

	validDuration, err := cert.ParseValidFor()
	if err != nil {
		return ctrl.Result{}, err
	}

	if time.Until(cert.Status.ExpiryDate.Time) <= 48*time.Hour {
		r.updateStatus(ctx, &cert, certv1.ConditionRenewing, "Certificate is being renewed")
		// Perfrom Renewal
		cert.Status.ExpiryDate = metav1.NewTime(time.Now().Add(validDuration))
		r.updateStatus(ctx, &cert, certv1.ConditionIssued, "Certificate successfully renewed")
	} else if time.Until(cert.Status.ExpiryDate.Time) <= 0 {
		r.updateStatus(ctx, &cert, certv1.ConditionExpired, "Certificate is expired")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&certv1.Certificate{}).
		Complete(r)
}
