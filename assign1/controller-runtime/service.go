package crt

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateService ...
func (c *CRTClient) CreateService(serviceName string) error {
	pod := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: c.namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       8091,
					TargetPort: intstr.FromInt(8091),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
	return c.apiServerClient.Create(context.Background(), pod, &client.CreateOptions{})
}

// UpdateService ...
func (c *CRTClient) UpdateService(serviceName string) error {
	oldService, err := c.GetService(serviceName)
	if err != nil {
		log.Println("failed to get old service : ", err)
		return err
	}

	oldService.SetGenerateName("updated-ctr-service")
	err = c.apiServerClient.Update(context.Background(), oldService, &client.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// DeleteService ...
func (c *CRTClient) DeleteService(serviceName string) error {
	existingService, err := c.GetService(serviceName)
	if err != nil {
		fmt.Println("error getting a service : ", err)
		return err
	}

	return c.apiServerClient.Delete(context.Background(), existingService)
}

// GetService ...
func (c *CRTClient) GetService(serviceName string) (*corev1.Service, error) {
	retrivedService := &corev1.Service{}
	err := c.apiServerClient.Get(context.Background(), types.NamespacedName{Name: serviceName, Namespace: c.namespace}, retrivedService)
	if err != nil {
		fmt.Println("error getting a service : ", err)
		return nil, err
	}

	return retrivedService, nil
}
