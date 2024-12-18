package internal

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type IK8sOperation interface {
	ListNamespaces(ctx context.Context) ([]*corev1.Namespace, error)

	ListDeployments(ctx context.Context, namespace string) ([]*appsv1.Deployment, error)
	GetDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error)
	DeleteDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error)
	CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
	UpdateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error)
}
