package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	dbv1 "velotio.com/database/api/v1"
)

func (r *MongoDBReconciler) upsertMongoPod(ctx context.Context, db *dbv1.MongoDB) error {

	d := &appsv1.Deployment{}
	// get deployment if exists
	if err := r.Get(ctx, types.NamespacedName{Namespace: db.GetNamespace(), Name: "mongodb-deployment"}, d); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		// deployment not found, create new
		d = getMongoDBResource(db)
		if err = controllerutil.SetControllerReference(db, d, r.Scheme); err != nil {
			return err
		}
		if err = r.Create(ctx, d); err != nil {
			return err
		}
		return nil
	}
	// pod exists, update the spec
	if err := controllerutil.SetControllerReference(db, d, r.Scheme); err != nil {
		return err
	}
	if err := r.Update(ctx, d); err != nil {
		return err
	}
	return nil
}

func getMongoDBResource(db *dbv1.MongoDB) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb-deployment",
			Namespace: db.GetNamespace(),
			Labels: map[string]string{
				"app": "mongodb",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mongodb",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "mongodb",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Image:           "mongo",
							ImagePullPolicy: "IfNotPresent",
							Name:            "mongodb",
							Env: []v1.EnvVar{
								{
									Name:  "MONGO_INITDB_ROOT_USERNAME",
									Value: db.Spec.InitUser,
								},
								{
									Name:  "MONGO_INITDB_ROOT_PASSWORD",
									Value: db.Spec.InitUser,
								},
							},
						},
						{
							Image:           "bitnami/kubectl",
							ImagePullPolicy: "IfNotPresent",
							Name:            "kubectl",
							Command: []string{
								"kubectl",
								"-n",
								"ass-4",
								"create",
								"job",
								"dummy-job",
								"--image",
								"busybox",
							},
						},
					},
					ServiceAccountName: "ass-4-sa",
				},
			},
		},
	}
}
