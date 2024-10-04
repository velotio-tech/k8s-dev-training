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
	"fmt"
	"github.com/jshiwamv/k8s-dev-training/certificate-operator/internal/certificate"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"
	"strings"
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
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	validDuration, err := parseValidFor(cert.Spec.ValidFor)
	if err != nil {
		return ctrl.Result{}, err
	}

	if time.Until(cert.Status.ExpiryDate.Time) <= 48*time.Hour {
		err = r.updateStatus(ctx, cert, certv1.ConditionRenewing, "Certificate is being renewed")
		if err != nil {
			return ctrl.Result{}, err
		}

		// Perfrom Renewal
		newCertPems, err := certificate.GenerateSelfSignedCertificate(validDuration)
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.createOrUpdateCertificate(ctx, cert, newCertPems)
		if err != nil {
			return ctrl.Result{}, err
		}

		cert.Status.ExpiryDate = metav1.NewTime(time.Now().Add(validDuration))

		err = r.updateStatus(ctx, cert, certv1.ConditionIssued, "Certificate successfully renewed")
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if time.Until(cert.Status.ExpiryDate.Time) <= 0 {
		err = r.updateStatus(ctx, cert, certv1.ConditionExpired, "Certificate is expired")
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
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

func (r *CertificateReconciler) updateStatus(ctx context.Context, cert *certv1.Certificate, condition certv1.CertificateConditionType, message string) error {
	newCondition := certv1.CertificateCondition{
		Type:               condition,
		Message:            message,
		LastTransitionTime: metav1.Now(),
		Status:             metav1.ConditionTrue,
	}
	cert.Status.Conditions = append(cert.Status.Conditions, newCondition)
	return r.Status().Update(ctx, cert)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CertificateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&certv1.Certificate{}).
		Complete(r)
}

func parseValidFor(validFor string) (time.Duration, error) {
	if strings.HasSuffix(validFor, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(validFor, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	if strings.HasSuffix(validFor, "y") {
		year, err := strconv.Atoi(strings.TrimSuffix(validFor, "y"))
		if err != nil {
			return 0, err
		}
		return time.Duration(year) * 365 * 24 * time.Hour, nil
	}
	return 0, fmt.Errorf("invalid value %s for validFor for field, should end with `d`(days) or `y`(years) e:g 1y, 20d ", validFor)
}
