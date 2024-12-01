package services

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

var services = &apiv1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Name: "demo-service",
	},
	Spec: apiv1.ServiceSpec{
		Selector: map[string]string{
			"app": "demo",
		},
		Ports: []apiv1.ServicePort{
			{
				Name:     "access-port",
				Protocol: "TCP",
				Port:     8080,
			},
		},
	},
}
var servicesClient corev1.ServiceInterface

func CreateServicesClient(clientset *kubernetes.Clientset) {
	servicesClient = clientset.CoreV1().Services(apiv1.NamespaceDefault)
}
func CreateServices() {
	// Create Service
	fmt.Println("Creating service...")
	result, err := servicesClient.Create(context.Background(), services, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	if err != nil {
		log.Println("Error occcured while creating service", err.Error())
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())
}

func ListServices() {
	// List Services
	fmt.Printf("Listing services in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := servicesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}

}

func UpdateServices() {
	// Update services
	//prompt()
	fmt.Println("Updating services...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of service before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := servicesClient.Get(context.Background(), "demo-service", metav1.GetOptions{})
		if getErr != nil {
			//panic(fmt.Errorf("Failed to get latest version of service: %v", getErr))
			log.Println(fmt.Errorf("Failed to get latest version of service: %v", getErr))
		}

		result.Spec.Ports[0].Protocol = "UDP" // update protocol to UDP
		_, updateErr := servicesClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated service...")
}

func DeleteServices() {
	// Delete Service
	//prompt()
	fmt.Println("Deleting service...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := servicesClient.Delete(context.Background(), "demo-service", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service.")

}
