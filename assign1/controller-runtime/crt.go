package crt

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// CRT ...
type CRT interface {
	CreatePod(podName string) error
	UpdatePod(podName string) error
	DeletePod(podName string) error
	GetPod(podName string) (*appsv1.Deployment, error)

	CreateService(serviceName string) error
	UpdateService(serviceName string) error
	DeleteService(serviceName string) error
	GetService(serviceName string) (*corev1.Service, error)
}

// CRTClient ...
type CRTClient struct {
	apiServerClient client.Client
	namespace       string
}

// NewCRTClient ...
func NewCRTClient(namespace string) (*CRTClient, error) {
	apiServerClient, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the crt client : %v", err)
	}
	return &CRTClient{
		apiServerClient: apiServerClient,
		namespace:       namespace,
	}, nil
}
