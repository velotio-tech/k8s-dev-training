package clientgo

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

var service = &apiv1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Name: "test-service",
	},
	Spec: apiv1.ServiceSpec{
		Selector: map[string]string{
			"app": "demo",
		},
		Ports: []apiv1.ServicePort{
			{
				Name:     "access-port",
				Protocol: "TCP",
				Port:     8009,
			},
		},
	},
}

func GetAllServices(serviceClient corev1.ServiceInterface) error {

	fmt.Printf("Listing services in namespace %s:\n", apiv1.NamespaceDefault)
	list, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}
	return nil
}

func CreateServices(serviceClient corev1.ServiceInterface) error {

	fmt.Println("Creating Service...")
	result, err := serviceClient.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created service %s.\n", result.GetObjectMeta().GetName())
	return nil
}

func DeleteService(serviceClient corev1.ServiceInterface) error {

	fmt.Println("Deleting service...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := serviceClient.Delete(context.Background(), "test-service", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	fmt.Println("Deleted service.")
	return nil
}

func UpdateServices(serviceClient corev1.ServiceInterface) error {

	fmt.Println("Updating service...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := serviceClient.Get(context.Background(), "test-service", metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("failed to get latest version of service: %v ", getErr)
		}

		result.Spec.Ports[0].Protocol = "UDP"
		_, updateErr := serviceClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return (fmt.Errorf("update failed: %v ", retryErr))
	}
	fmt.Println("Updated service...")
	return nil
}
