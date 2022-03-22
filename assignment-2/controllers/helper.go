package controllers

import (
	"context"
	v1 "github.com/hatred09/k8s-dev-training/assignment-2/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SecretCreatorReconciler) CreateSecret(ctx context.Context, namespace string, secretcreator v1.SecretCreator, data map[string][]byte) error {
	secretDeploy := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SecretCreator",
			APIVersion: "secretcreator.example.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      secretcreator.Spec.SecretName,
		},
		Data: data,
	}
	ctrl.SetControllerReference(&secretcreator, secretDeploy, r.Scheme)
	err := r.Create(ctx, secretDeploy)
	if err != nil {
		return err
	}
	return nil
}

func (r *SecretCreatorReconciler) DeleteSecret(ctx context.Context, namespace, name string, data map[string][]byte) error {
	secretDeploy := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Data: data,
	}
	err := r.Delete(ctx, secretDeploy)
	if err != nil {
		return err
	}
	return nil
}

func (r *SecretCreatorReconciler) GetNamespaces(ctx context.Context) (*corev1.NamespaceList, error) {
	namespaces := &corev1.NamespaceList{}
	if err := r.List(ctx, namespaces); err != nil {
		return nil, err
	}
	return namespaces, nil
}

func contains(slice []string, item string) bool {
	set := make(sets.String)
	for _, s := range slice {
		set[s] = struct{}{}
	}
	return set.Has(item)
}
