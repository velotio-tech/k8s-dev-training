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
	"github.com/jshiwamv/k8s-dev-training/certificate-operator/internal/certificate"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
	log := log.FromContext(ctx)
	cert := &certv1.Certificate{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, cert)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "Certificate resource not found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to fetch certificate")
		return ctrl.Result{}, err
	}

	validDuration, err := cert.Spec.ParseValidFor()
	if err != nil {
		log.Error(err, "Failed to parse ValidFor Field from spec")
		return ctrl.Result{}, err
	}
	log.Info("Reconcilation Triggered ....")

	if cert.Status.ExpiryDate.IsZero() {
		log.Info("Creating new Cert ....")
		err = r.updateStatus(ctx, req, metav1.Condition{
			Type:    certv1.ConditionPending.String(),
			Status:  metav1.ConditionTrue,
			Reason:  "Reconciling",
			Message: "Certificate creation is pending",
		}, nil, nil)
		if err != nil {
			log.Error(err, "Failed to update certificate status", "ConditionType", certv1.ConditionPending.String())
			return ctrl.Result{}, err
		}

		// Perfrom Renewal
		newCertPems, err := certificate.GenerateSelfSignedCertificate(validDuration, cert.Spec.Domain)
		if err != nil {
			log.Error(err, "Failed to generate self signed certificate pem values for certificate creation")
			return ctrl.Result{}, err
		}

		err = r.createOrUpdateCertificate(ctx, cert, newCertPems)
		if err != nil {
			log.Error(err, "Failed to create certificate")
			return ctrl.Result{}, err
		}

		expiredAt := metav1.NewTime(time.Now().Add(validDuration))
		err = r.updateStatus(ctx, req, metav1.Condition{
			Type:    certv1.ConditionIssued.String(),
			Status:  metav1.ConditionTrue,
			Reason:  "Reconciling",
			Message: "Certificate successfully issued",
		}, &expiredAt, nil)
		if err != nil {
			log.Error(err, "Failed to update certificate status", "ConditionType", certv1.ConditionIssued.String())
			return ctrl.Result{}, err
		}
	} else if time.Until(cert.Status.ExpiryDate.Time) <= 5*time.Minute {
		log.Info("Renewing expired cert ....")
		err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, cert)
		if err != nil {
			log.Error(err, "Failed to re-fetch certificate")
			return ctrl.Result{}, err
		}

		err = r.updateStatus(ctx, req, metav1.Condition{
			Type:    certv1.ConditionRenewing.String(),
			Status:  metav1.ConditionTrue,
			Reason:  "Reconciling",
			Message: "Certificate renewal in progress",
		}, nil, nil)
		if err != nil {
			log.Error(err, "Failed to update certificate status", "ConditionType", certv1.ConditionRenewing.String())
			return ctrl.Result{}, err
		}

		// Perfrom Renewal
		newCertPems, err := certificate.GenerateSelfSignedCertificate(validDuration, cert.Spec.Domain)
		if err != nil {
			log.Error(err, "Failed to generate self signed certificate pem values for certificate renewal")
			return ctrl.Result{}, err
		}

		err = r.createOrUpdateCertificate(ctx, cert, newCertPems)
		if err != nil {
			log.Error(err, "Failed to renew certificate")
			return ctrl.Result{}, err
		}

		renewedAt := metav1.Now()
		expiredAt := metav1.NewTime(time.Now().Add(validDuration))

		err = r.updateStatus(ctx, req, metav1.Condition{
			Type:    certv1.ConditionRenewed.String(),
			Status:  metav1.ConditionTrue,
			Reason:  "Reconciling",
			Message: "Certificate successfully renewed",
		}, &expiredAt, &renewedAt)
		if err != nil {
			log.Error(err, "Failed to update certificate status", "ConditionType", certv1.ConditionRenewed.String())
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: time.Minute * 1}, nil
}

func (r *CertificateReconciler) createOrUpdateCertificate(ctx context.Context, cert *certv1.Certificate, certPEM *certificate.CertificatePEM) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cert.Namespace,
			Name:      cert.Name,
		},
		Data: map[string][]byte{
			"tls.crt": certPEM.CertPEM,
			"tls.key": certPEM.KeyPEM,
		},
	}

	if err := controllerutil.SetOwnerReference(cert, secret, r.Scheme); err != nil {
		return err
	}

	err := r.Client.Patch(ctx, secret, client.MergeFrom(secret))
	if err != nil {
		return r.Client.Create(ctx, secret)
	}
	return nil
}

func (r *CertificateReconciler) updateStatus(ctx context.Context, req ctrl.Request, condition metav1.Condition, expiredAt, renewedAt *metav1.Time) error {
	cert := &certv1.Certificate{}

	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, cert)
	if err != nil {
		return err
	}

	newCondition := metav1.Condition{
		Type:               condition.Type,
		Message:            condition.Message,
		LastTransitionTime: metav1.Now(),
		Status:             condition.Status,
		Reason:             condition.Reason,
	}

	for i := range cert.Status.Conditions {
		cond := &cert.Status.Conditions[i]
		if cond.Type != newCondition.Type {
			cond.Status = metav1.ConditionFalse
		}
	}

	if expiredAt != nil {
		cert.Status.ExpiryDate = *expiredAt
	}
	if renewedAt != nil {
		cert.Status.RenewedAt = *renewedAt
	}

	meta.SetStatusCondition(&cert.Status.Conditions, newCondition)
	return r.Status().Update(ctx, cert)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&certv1.Certificate{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
