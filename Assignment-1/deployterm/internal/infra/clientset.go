package infra

import (
	"context"

	"github.com/aryanmaurya1/deployterm/internal/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type clientsetWrapper struct {
	clientset *kubernetes.Clientset
}

func NewClientsetWrapper(clientset *kubernetes.Clientset) *clientsetWrapper {
	return &clientsetWrapper{clientset: clientset}
}

func (csw *clientsetWrapper) ListNamespaces(ctx context.Context) ([]*corev1.Namespace, error) {
	namespaceList, err := csw.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return utils.ConvertToPtrList[corev1.Namespace](namespaceList.Items), nil
}

func (csw *clientsetWrapper) ListDeployments(ctx context.Context, namespace string) ([]*appsv1.Deployment, error) {
	list, err := csw.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return utils.ConvertToPtrList(list.Items), nil
}

func (csw *clientsetWrapper) GetDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	return csw.clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
}

func (csw *clientsetWrapper) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return csw.clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
}

func (csw *clientsetWrapper) UpdateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return csw.clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
}

func (csw *clientsetWrapper) DeleteDeployment(ctx context.Context, namespace string, deploymentName string) (*appsv1.Deployment, error) {
	deployment, err := csw.GetDeployment(ctx, namespace, deploymentName)
	if err != nil {
		return nil, err
	}

	err = csw.clientset.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, err
}
