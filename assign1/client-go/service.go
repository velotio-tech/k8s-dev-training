package clientgo

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateService ...
func (cg *ClientGoClient) CreateService() (*corev1.Service, error) {
	newService := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cg.setUpEssentials.ServiceName,
			Namespace: cg.setUpEssentials.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "arijitAssign1",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       8081,
					TargetPort: intstr.FromInt(8081),
				},
			},
		},
	}

	createdService, err := cg.serviceClient.Create(context.Background(), newService, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createdService, nil
}

// UpdateService ...
func (cg *ClientGoClient) UpdateService(serviceName string) error {
	updatedService, err := cg.serviceClient.Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting Pod for update: %v", err)
	}

	updatedService.SetGenerateName("generated-updated-service")

	_, err = cg.serviceClient.Update(context.Background(), updatedService, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating Pod: %v", err)
	}

	return nil
}

// GetServices ...
func (cg *ClientGoClient) GetServices() (*corev1.ServiceList, error) {
	return cg.serviceClient.List(context.Background(), metav1.ListOptions{})
}

// DeleteService ...
func (cg *ClientGoClient) DeleteService(podName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return cg.serviceClient.Delete(context.Background(), podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
