package controllers

import (
	"context"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	dbv1 "velotio.com/database/api/v1"
)

func (r *MongoDBReconciler) upsertMongoPod(ctx context.Context, db *dbv1.MongoDB) error {

	p := &v1.Pod{}
	// get pod if exists
	if err := r.Get(ctx, types.NamespacedName{Namespace: db.GetNamespace(), Name: "mongodb-pod"}, p); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		// pod not found, create new
		pod := getMongoDBResource(db)
		if err = controllerutil.SetControllerReference(db, pod, r.Scheme); err != nil {
			return err
		}
		if err = r.Create(ctx, pod); err != nil {
			return err
		}
		return nil
	}
	// pod exists, update the spec
	if err := controllerutil.SetControllerReference(db, p, r.Scheme); err != nil {
		return err
	}
	if err := r.Update(ctx, p); err != nil {
		return err
	}
	return nil
}

func getMongoDBResource(db *dbv1.MongoDB) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb-pod",
			Namespace: db.GetNamespace(),
			Labels: map[string]string{
				"app": "mongodb",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "mongodb",
					Image: "mongo:latest",
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
			},
		},
	}
}
