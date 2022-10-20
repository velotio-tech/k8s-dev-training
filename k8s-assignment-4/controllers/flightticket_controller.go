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
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

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

	flightTicket := &k8sassignment2v1.FlightTicket{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, flightTicket)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Unable to fetch FlightTicket - skipping")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch FlightTicket")
		return ctrl.Result{}, err
	}

	err = r.FlightTicketGvk(flightTicket, log)
	if err != nil {
		log.Error(err, "Failed to create/update bookstore resources")
		return ctrl.Result{}, err
	}

	finalizerName := "k8s-assignment-2/flight-ticket/finalizer"
	// Since, delete on an object is Soft-delete, the presence of deletion timestamp on the object indicates that it is being deleted.
	if flightTicket.ObjectMeta.DeletionTimestamp.IsZero() {
		// If the object is not being deleted and does not have the finalizer registered,
		// then add the finalizer and update the object in Kubernetes.
		if !controllerutil.ContainsFinalizer(flightTicket, finalizerName) {
			controllerutil.AddFinalizer(flightTicket, finalizerName)
			if err := r.Update(ctx, flightTicket); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// If object is being deleted and the finalizer is still present in finalizers list,
		// then execute the pre-delete logic and remove the finalizer and update the object.
		if controllerutil.ContainsFinalizer(flightTicket, finalizerName) {
			if err := r.cleanupExternalResources(flightTicket, log); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(flightTicket, finalizerName)
			if err := r.Update(ctx, flightTicket); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	flightTicket.Status.ReadyReplicas = int(flightTicket.Spec.Gvk.Replicas)

	if flightTicket.Status.BookingStatus != "Done" {
		flightTicket.Status.BookingStatus = "Done"
		flightTicket.Status.Fare = 5000
		if err := r.Status().Update(ctx, flightTicket); err != nil {
			log.Error(err, "unable to update flightticket status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FlightTicketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	indexerFunc := func(rawObj client.Object) []string {
		flightTicket := rawObj.(*k8sassignment2v1.FlightTicket)
		owner := metav1.GetControllerOf(flightTicket)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != k8sassignment2v1.GroupVersion.String() || owner.Kind != "flightTicket" {
			return nil
		}
		return []string{owner.Name}
	}

	err := mgr.GetFieldIndexer().IndexField(context.Background(), &k8sassignment2v1.FlightTicket{}, ".metadata.controller", indexerFunc)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sassignment2v1.FlightTicket{}).
		WithEventFilter(r.ignoreDeletionPredicate()).
		Complete(r)
}

func (r *FlightTicketReconciler) cleanupExternalResources(flightTicket *k8sassignment2v1.FlightTicket, log logr.Logger) error {
	// cancel booking and charge cancellation fee
	log.Info("cancel booking and charge cancellation fee")
	flightTicket.Status.BookingStatus = "Cancelled"
	flightTicket.Status.Fare = 350

	return nil
}

func (r *FlightTicketReconciler) FlightTicketGvk(flightTicket *k8sassignment2v1.FlightTicket, log logr.Logger) error {

	switch flightTicket.Spec.Gvk.Kind {
	case "StatefulSet":
		flightTicketStatefulSet := getFlightTicketStatefulSet(flightTicket)
		statefulSet := &appsv1.StatefulSet{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{Name: flightTicketStatefulSet.Name, Namespace: flightTicketStatefulSet.Namespace}, statefulSet)
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Info("FlightTicket StatefulSet not found, will be created")
				controllerutil.SetControllerReference(flightTicket, flightTicketStatefulSet, r.Scheme)
				err = r.Client.Create(context.TODO(), flightTicketStatefulSet)
				if err != nil {
					return err
				}
			} else {
				log.Info("failed to get FlightTicket StatefulSet")
				return err
			}
		} else if !reflect.DeepEqual(flightTicketStatefulSet.Spec, statefulSet.Spec) {
			flightTicketStatefulSet.ObjectMeta = statefulSet.ObjectMeta

			controllerutil.SetControllerReference(flightTicket, flightTicketStatefulSet, r.Scheme)
			err = r.Client.Update(context.TODO(), flightTicketStatefulSet)
			if err != nil {
				return err
			}
			log.Info("FlightTicket StatefulSet updated")
		}
	case "Deployment":
		flightTicketDeployment := getFlightTicketDeployment(flightTicket)
		dss := &appsv1.Deployment{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{Name: flightTicketDeployment.Name, Namespace: flightTicketDeployment.Namespace}, dss)
		if err != nil {
			if apierrors.IsNotFound(err) {
				log.Info("FlightTicket Deployment not found, will be created")
				controllerutil.SetControllerReference(flightTicket, flightTicketDeployment, r.Scheme)
				err = r.Client.Create(context.TODO(), flightTicketDeployment)
				if err != nil {
					return err
				}
			} else {
				log.Info("failed to get FlightTicket Deployment")
				return err
			}
		} else if !reflect.DeepEqual(flightTicketDeployment.Spec, dss.Spec) {
			flightTicketDeployment.ObjectMeta = dss.ObjectMeta
			controllerutil.SetControllerReference(flightTicket, flightTicketDeployment, r.Scheme)
			err = r.Client.Update(context.TODO(), flightTicketDeployment)
			if err != nil {
				return err
			}
			log.Info("FlightTicket Deployment updated")
		}
	default:
		str := fmt.Sprintf("Invalid value for kind (must be either StatefulSet or Deployment): %v\n", flightTicket.Spec.Gvk.Kind)
		log.Error(nil, str)
	}

	return nil
}

func getFlightTicketStatefulSet(flightTicket *k8sassignment2v1.FlightTicket) *appsv1.StatefulSet {
	labelMap := map[string]string{"app": "flight-ticket"}
	containerArr := make([]corev1.Container, 0)
	container := corev1.Container{
		Name:  "flight-ticket",
		Image: "nginx",
	}
	containerArr = append(containerArr, container)

	podTempSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labelMap,
		},
		Spec: corev1.PodSpec{
			Containers: containerArr,
		},
	}
	flightTicketStatefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flight-ticket",
			Namespace: flightTicket.Namespace,
			Labels:    labelMap,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelMap,
			},
			Replicas:    &flightTicket.Spec.Gvk.Replicas,
			Template:    podTempSpec,
			ServiceName: "flight-ticket-service",
		},
	}
	return flightTicketStatefulSet
}

func getFlightTicketDeployment(flightTicket *k8sassignment2v1.FlightTicket) *appsv1.Deployment {
	labelMap := map[string]string{"app": "flight-ticket"}
	containerArr := make([]corev1.Container, 0)
	container := corev1.Container{
		Name:  "flight-ticket",
		Image: "nginx",
	}
	containerArr = append(containerArr, container)

	podTempSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labelMap,
		},
		Spec: corev1.PodSpec{
			Containers: containerArr,
		},
	}
	flightTicketDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      flightTicket.Name,
			Namespace: flightTicket.Namespace,
			Labels:    labelMap,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelMap,
			},
			Replicas: &flightTicket.Spec.Gvk.Replicas,
			Template: podTempSpec,
		},
	}
	return flightTicketDeployment
}

func (r *FlightTicketReconciler) ignoreDeletionPredicate() predicate.Predicate {
	ctx := context.Background()
	log := log.FromContext(ctx)
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			log.Info("In update func predicate")
			updatedTicket, ok := e.ObjectNew.(*k8sassignment2v1.FlightTicket)
			if ok {
				if updatedTicket.Status.BookingStatus == "InProgress" {
					return false
				}

				oldTicket, _ := e.ObjectOld.(*k8sassignment2v1.FlightTicket)
				if reflect.DeepEqual(updatedTicket.Spec, oldTicket.Spec) && oldTicket.Status.BookingStatus != "" {
					return false
				}
			}

			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			log.Info("In delete func predicate")

			if obj, ok := e.Object.(*k8sassignment2v1.FlightTicket); ok {
				if obj.Status.BookingStatus != "InProgress" {
					return true
				}
			}
			return false
		},
	}
}
