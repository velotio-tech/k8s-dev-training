package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appsv1alpha1 "github.com/Ajay-Raut-Rsystems/k8s-dev-training/assignment4/api/v1alpha1"
)

const (
	ParentAnnotationKey = "apps.example.com/parent-cr"
	JobFieldIndexerKey  = "metadata.annotations.parent-cr"
)

type MultiTierAppReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	WatchNamespace string // Target namespace to enforce
}

// +kubebuilder:rbac:groups=apps.example.com,resources=multitierapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.example.com,resources=multitierapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

func (r *MultiTierAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. Fetch Level 1 Parent CR
	parent := &appsv1alpha1.MultiTierApp{}
	if err := r.Get(ctx, req.NamespacedName, parent); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// =========================================================
	// LEVEL 2: Reconcile Deployment (WITH OwnerReference)
	// =========================================================
	deployName := fmt.Sprintf("%s-deploy", parent.Name)
	foundDeploy := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: deployName, Namespace: parent.Namespace}, foundDeploy)

	if errors.IsNotFound(err) {
		deploy := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployName,
				Namespace: parent.Namespace,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &parent.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": deployName},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": deployName},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:    "app-container",
							Image:   "busybox",
							Command: []string{"sh", "-c", "echo Level 2 running && sleep 3600"},
						}},
					},
				},
			},
		}

		// Set Level 1 -> Level 2 OwnerReference
		if err := controllerutil.SetControllerReference(parent, deploy, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		logger.Info("Creating Level 2 Deployment", "Deployment", deploy.Name)
		if err := r.Create(ctx, deploy); err != nil {
			return ctrl.Result{}, err
		}
	}

	// =========================================================
	// LEVEL 3: Reconcile Job (WITHOUT OwnerReference, using Field Indexer)
	// =========================================================
	var jobList batchv1.JobList
	// Query cache via Field Indexer key
	err = r.List(ctx, &jobList,
		client.InNamespace(parent.Namespace),
		client.MatchingFields{JobFieldIndexerKey: parent.Name},
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(jobList.Items) == 0 {
		jobName := fmt.Sprintf("%s-job", parent.Name)
		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: parent.Namespace,
				Annotations: map[string]string{
					ParentAnnotationKey: parent.Name, // Annotation used for field indexing
				},
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyNever,
						Containers: []corev1.Container{{
							Name:    "job-task",
							Image:   "busybox",
							Command: []string{"sh", "-c", parent.Spec.JobCommand},
						}},
					},
				},
			},
		}

		// NOTE: Notice we DO NOT call controllerutil.SetControllerReference for Level 3!
		logger.Info("Creating Level 3 Job (Unowned)", "Job", job.Name)
		if err := r.Create(ctx, job); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 4. Re-fetch latest instance before updating status to avoid Optimistic Concurrency Conflicts
	latestParent := &appsv1alpha1.MultiTierApp{}
	if err := r.Get(ctx, req.NamespacedName, latestParent); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Update Status
	latestParent.Status.Phase = "Completed"
	latestParent.Status.DeploymentName = deployName
	latestParent.Status.JobName = fmt.Sprintf("%s-job", parent.Name)
	if err := r.Status().Update(ctx, latestParent); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// Namespace-Scoped Predicate
func (r *MultiTierAppReconciler) namespacePredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			if r.WatchNamespace == "" {
				return true
			}
			return e.Object.GetNamespace() == r.WatchNamespace
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			if r.WatchNamespace == "" {
				return true
			}
			return e.ObjectNew.GetNamespace() == r.WatchNamespace
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			if r.WatchNamespace == "" {
				return true
			}
			return e.Object.GetNamespace() == r.WatchNamespace
		},
	}
}

func (r *MultiTierAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Register Field Indexer for Job Annotations
	err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&batchv1.Job{},
		JobFieldIndexerKey,
		func(rawObj client.Object) []string {
			job := rawObj.(*batchv1.Job)
			if parentName, exists := job.Annotations[ParentAnnotationKey]; exists {
				return []string{parentName}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.MultiTierApp{}).
		WithEventFilter(r.namespacePredicate()). // Apply Namespace Scope
		Owns(&appsv1.Deployment{}).              // Watch Level 2 Owned Deployment
		Complete(r)
}
