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
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	myappsv1 "assignment.com/dummydb/api/v1"
)

var log = logf.Log.WithName("controller_dummydb")

// DummyDBReconciler reconciles a DummyDB object
type DummyDBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapps.assignment.com,resources=dummydbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapps.assignment.com,resources=dummydbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapps.assignment.com,resources=dummydbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DummyDB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *DummyDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

	dummyDB := &myappsv1.DummyDB{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, dummyDB)
	if err != nil {
		if errors.IsNotFound(err) {

			err = r.CleanUpPVCAfterDBRemoval(ctx, req)

			if err != nil {
				reqLogger.Error(err, "Failed to delete PVCs")
			}

			return reconcile.Result{}, err
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get dummydb object")
		return reconcile.Result{}, err
	}

	dummydbFinalizer := "assignment.dummydb.io/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if dummyDB.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(dummyDB, dummydbFinalizer) {
			controllerutil.AddFinalizer(dummyDB, dummydbFinalizer)
			if err := r.Update(ctx, dummyDB); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(dummyDB, dummydbFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.CleanUpPVCAfterDBRemoval(ctx, req); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(dummyDB, dummydbFinalizer)
			if err := r.Update(ctx, dummyDB); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	err = r.DummyDB(dummyDB)
	if err != nil {
		reqLogger.Error(err, "Failed to create/update bookstore resources")
		return reconcile.Result{}, err
	}

	// update status
	dummyDB.Status.AvailabeSize = int32(dummyDB.Spec.Size.AsApproximateFloat64())

	dummyDB.Status.ReadyReplicas = dummyDB.Spec.Replicas

	dummyDB.Status.VolumeResizingInProgress = false

	_ = r.Client.Status().Update(context.TODO(), dummyDB)

	return ctrl.Result{}, nil
}

func (r *DummyDBReconciler) DummyDB(dummyDB *myappsv1.DummyDB) error {

	reqLogger := log.WithValues("Namespace", dummyDB.Namespace)

	if dummyDB.Spec.Gvk.Kind == "StatefulSet" {
		dummyDBSS := getDummyDBStatefulset(dummyDB)
		dss := &appsv1.StatefulSet{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{Name: dummyDBSS.Name, Namespace: dummyDBSS.Namespace}, dss)
		if err != nil {
			if errors.IsNotFound(err) {
				reqLogger.Info("dummyDB statefulset not found, will be created")
				controllerutil.SetControllerReference(dummyDB, dummyDBSS, r.Scheme)
				err = r.Client.Create(context.TODO(), dummyDBSS)
				if err != nil {
					return err
				}
			} else {
				reqLogger.Info("failed to get dummyDB statefulset")
				return err
			}
		} else if !reflect.DeepEqual(dummyDBSS.Spec, dss.Spec) {
			//r.UpdateVolume(dummyDB)
			dummyDBSS.ObjectMeta = dss.ObjectMeta
			dummyDBSS.Spec.VolumeClaimTemplates = dss.Spec.VolumeClaimTemplates
			controllerutil.SetControllerReference(dummyDB, dummyDBSS, r.Scheme)
			err = r.Client.Update(context.TODO(), dummyDBSS)
			if err != nil {
				return err
			}
			reqLogger.Info("dummyDB statefulset updated")
		}

	} else if dummyDB.Spec.Gvk.Kind == "Deployment" {

		dummyDBDeployment := getDummyDBDeploy(dummyDB)
		dss := &appsv1.Deployment{}
		err := r.Client.Get(context.TODO(), types.NamespacedName{Name: dummyDBDeployment.Name, Namespace: dummyDBDeployment.Namespace}, dss)
		if err != nil {
			if errors.IsNotFound(err) {
				reqLogger.Info("dummyDB Deployment not found, will be created")
				volClaimTemplate(dummyDB.Spec.Size)
				controllerutil.SetControllerReference(dummyDB, dummyDBDeployment, r.Scheme)
				err = r.Client.Create(context.TODO(), dummyDBDeployment)
				if err != nil {
					return err
				}
			} else {
				reqLogger.Info("failed to get dummyDB Deployment")
				return err
			}
		} else if !reflect.DeepEqual(dummyDBDeployment.Spec, dss.Spec) {
			//r.UpdateVolume(dummyDB)
			dummyDBDeployment.ObjectMeta = dss.ObjectMeta
			controllerutil.SetControllerReference(dummyDB, dummyDBDeployment, r.Scheme)
			err = r.Client.Update(context.TODO(), dummyDBDeployment)
			if err != nil {
				return err
			}
			reqLogger.Info("dummyDB Deployment updated")
		}

	} else {
		reqLogger.Error(nil, "Invalid value for kind (must be StatefulSet or Deployment): %v\n", dummyDB.Spec.Gvk.Kind)
	}

	dummyDBSecret := getDummyDBSecret(dummyDB)
	dds := &v1.Secret{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: dummyDBSecret.Name, Namespace: dummyDBSecret.Namespace}, dds)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("dummyDB's admin password secret not found, will be created")
			controllerutil.SetControllerReference(dummyDB, dummyDBSecret, r.Scheme)
			err = r.Client.Create(context.TODO(), dummyDBSecret)
			if err != nil {
				return err
			}
		} else {
			reqLogger.Info("failed to get dummyDB statefulset")
			return err
		}
	} else if !reflect.DeepEqual(dummyDBSecret, dds) {
		//r.UpdateVolume(dummyDB)
		dummyDBSecret.ObjectMeta = dds.ObjectMeta
		controllerutil.SetControllerReference(dummyDB, dummyDBSecret, r.Scheme)
		err = r.Client.Update(context.TODO(), dummyDBSecret)
		if err != nil {
			return err
		}
		reqLogger.Info("dummyDB's admin password secret updated")
	}

	return nil

}

func getDummyDBStatefulset(dummyDB *myappsv1.DummyDB) *appsv1.StatefulSet {

	cnts := make([]corev1.Container, 0)
	cnt := corev1.Container{
		Name:  "dummydb",
		Image: "nginx",
	}
	cnts = append(cnts, cnt)
	podTempSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"app": "dummydb"},
		},
		Spec: corev1.PodSpec{
			Containers: cnts,
		},
	}
	dummydb := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dummydb",
			Namespace: dummyDB.Namespace,
			Labels:    map[string]string{"app": "dummydb"},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "dummydb"},
			},
			Replicas:             &dummyDB.Spec.Replicas,
			Template:             podTempSpec,
			ServiceName:          "dummydb-service",
			VolumeClaimTemplates: volClaimTemplate(dummyDB.Spec.Size),
		},
	}
	return dummydb
}

func volClaimTemplate(dbSize resource.Quantity) []corev1.PersistentVolumeClaim {

	storageClass := "standard"
	resourceRequirements := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			//corev1.ResourceStorage: resource.MustParse(DBSize),
			corev1.ResourceStorage: dbSize,
		},
	}
	accessModeList := make([]corev1.PersistentVolumeAccessMode, 0)
	accessModeList = append(accessModeList, corev1.ReadWriteOnce)
	dummyDBPVC := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "dummydb-pvc",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      accessModeList,
			Resources:        resourceRequirements,
			StorageClassName: &storageClass,
		},
	}
	pvcList := make([]corev1.PersistentVolumeClaim, 0)
	pvcList = append(pvcList, dummyDBPVC)
	return pvcList
}

func getDummyDBDeploy(dummyDB *myappsv1.DummyDB) *appsv1.Deployment {

	cnts := make([]corev1.Container, 0)
	//pvcSource := v1.VolumeSource{PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: "dummydb-pvc"}}
	cnt := corev1.Container{
		Name:  "dummydb",
		Image: "nginx",
	}
	cnts = append(cnts, cnt)
	podTempSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"app": "dummydb"},
		},
		Spec: corev1.PodSpec{
			Containers: cnts,
			Volumes:    []v1.Volume{{Name: "dummydb"}},
		},
	}
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummyDB.Name,
			Namespace: dummyDB.Namespace,
			Labels:    map[string]string{"app": "dummydb"},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "dummydb"},
			},
			Replicas: &dummyDB.Spec.Replicas,
			Template: podTempSpec,
		},
	}
	return dep
}

func getDummyDBSecret(dummyDB *myappsv1.DummyDB) *v1.Secret {

	adminPasswordSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      dummyDB.Name,
			Namespace: dummyDB.Namespace,
			Labels:    map[string]string{"app": "dummydb"},
		},
		Data: map[string][]byte{
			"adminPassword": []byte(dummyDB.Spec.AdminPassword)},
	}

	return adminPasswordSecret
}

func (r *DummyDBReconciler) CleanUpPVCAfterDBRemoval(ctx context.Context, req ctrl.Request) error {

	reqLogger := log.WithValues("Namespace", req.Namespace, "Name", req.Name)

	pvcList := &corev1.PersistentVolumeClaimList{}

	if err := r.List(ctx, pvcList, client.InNamespace(req.Namespace), client.MatchingFields{OwnerKey: OwnerValue}); err != nil {
		reqLogger.Error(err, "unable to list pvcs")
		return err
	}

	for _, pvc := range pvcList.Items {
		err := r.Delete(ctx, &pvc)
		if err != nil {
			return err
		} else {
			reqLogger.Info("PVC deleted")
			return nil
		}
	}

	return nil
}

var (
	OwnerKey   = ".metadata.controller"
	OwnerValue = "DummyDB"
)

// SetupWithManager sets up the controller with the Manager.
func (r *DummyDBReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// PVC is created by statefulset and doesn't have resource owner referance set so setting an index for it.
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &corev1.PersistentVolumeClaim{}, OwnerKey, func(rawObj client.Object) []string {
		return []string{OwnerValue}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappsv1.DummyDB{}).WithEventFilter(onlyAllowSystemNS()).
		Complete(r)
}

func onlyAllowSystemNS() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			allowedNS := e.ObjectNew.GetNamespace()
			return allowedNS == "system"
		},
		CreateFunc: func(e event.CreateEvent) bool {
			allowedNS := e.Object.GetNamespace()
			return allowedNS == "system"
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			allowedNS := e.Object.GetNamespace()
			return allowedNS == "system"
		},
	}
}
