package infra

import (
	"context"

	"github.com/aryanmaurya1/deployterm/internal/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type runclientWrapper struct {
	runclient client.Client
}

func NewRunclientWrapper(runclient client.Client) *runclientWrapper {
	return &runclientWrapper{runclient: runclient}
}

func (rcw *runclientWrapper) ListNamespaces(ctx context.Context) ([]*corev1.Namespace, error) {
	var namespaceList corev1.NamespaceList
	err := rcw.runclient.List(ctx, &namespaceList)
	if err != nil {
		return nil, err
	}

	return utils.ConvertToPtrList[corev1.Namespace](namespaceList.Items), nil
}

func (rcw *runclientWrapper) ListDeployments(ctx context.Context, namespace string) ([]*appsv1.Deployment, error) {
	var deploymentList appsv1.DeploymentList
	err := rcw.runclient.List(ctx, &deploymentList, &client.ListOptions{Namespace: namespace})
	if err != nil {
		return nil, err
	}

	return utils.ConvertToPtrList[appsv1.Deployment](deploymentList.Items), nil
}

func (rcw *runclientWrapper) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	namespacedName := types.NamespacedName{Namespace: namespace, Name: deploymentName}
	var deployment appsv1.Deployment
	err := rcw.runclient.Get(ctx, namespacedName, &deployment)
	if err != nil {
		return nil, err
	}

	return &deployment, nil
}

func (rcw *runclientWrapper) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	deployment.Namespace = namespace
	err := rcw.runclient.Create(ctx, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (rcw *runclientWrapper) UpdateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	deployment.Namespace = namespace
	err := rcw.runclient.Update(ctx, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (rcw *runclientWrapper) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	deployment, err := rcw.GetDeployment(ctx, namespace, deploymentName)
	if err != nil {
		return nil, err
	}

	err = rcw.runclient.Delete(ctx, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}
