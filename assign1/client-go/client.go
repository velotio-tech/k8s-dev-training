package clientgo

import (
	"assig1/models"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	appV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

// ClientGo ...
type ClientGo interface {
	CreatePod() (*core.Pod, error)
	UpdatePod(podName string) error
	GetPods(podName string) (*core.PodList, error)
	DeletePod(podName string) error

	CreateDeployment() (*appsv1.Deployment, error)
	UpdateDeployment(deploymentName string) error
	GetDeployments() (*appsv1.DeploymentList, error)
	DeleteDeployment(deploymentName string) error

	CreateService() (*corev1.Service, error)
	UpdateService(serviceName string) error
	GetServices() (*corev1.ServiceList, error)
	DeleteService(podName string) error
}

// ClientGoClient ...
type ClientGoClient struct {
	podClient        v1.PodInterface
	serviceClient    v1.ServiceInterface
	deploymentClient appV1.DeploymentInterface
	setUpEssentials  *models.ClientGoEssentials
}

// NewClientGoClient ...
func NewClientGoClient(setUpEssentials *models.ClientGoEssentials) (*ClientGoClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Error creating in-cluster config: %v\n", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating clientset: %v\n", err)
		return nil, err
	}

	podClient := clientset.CoreV1().Pods(setUpEssentials.Namespace)
	serviceClient := clientset.CoreV1().Services(setUpEssentials.Namespace)
	deploymentClient := clientset.AppsV1().Deployments(setUpEssentials.Namespace)
	return &ClientGoClient{
		setUpEssentials:  setUpEssentials,
		podClient:        podClient,
		serviceClient:    serviceClient,
		deploymentClient: deploymentClient,
	}, nil
}
