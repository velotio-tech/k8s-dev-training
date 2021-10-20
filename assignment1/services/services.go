package services

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

var service = &apiv1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Name: "my-service",
	},
	Spec: apiv1.ServiceSpec{
		Selector: map[string]string{
			"app":"demo",
		},
		Ports: []apiv1.ServicePort{
			{
				Name: "access-port",
				Protocol: "TCP",
				Port: 8009,
			},
		},
	},
}

var serviceClient corev1.ServiceInterface


func CreateServicesClient(clientset *kubernetes.Clientset)  {
	serviceClient = clientset.CoreV1().Services(apiv1.NamespaceDefault)
}


func GetAllServices() {
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}
}

func CreateServices() {
	fmt.Println("Creating Service...")
	result, err := serviceClient.Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		log.Println("Error Occurred while creating the service", err.Error())
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
}

func DeleteService(){
	fmt.Println("Deleting service...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := serviceClient.Delete(context.Background(), "my-service", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service.")
}

func UpdateServices() {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := serviceClient.Get(context.Background(), "my-service", metav1.GetOptions{})
		if getErr != nil {
			log.Println(fmt.Errorf("Failed to get latest version of service: %v ", getErr))
		}

		result.Spec.Ports[0].Protocol = "UDP"
		_, updateErr := serviceClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v ", retryErr))
	}
	fmt.Println("Updated service...")
}