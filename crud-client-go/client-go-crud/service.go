package clientgocrud

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

func CreateService(serviceClient v1.ServiceInterface, clientset *kubernetes.Clientset) {

	//reference: https://stackoverflow.com/questions/53874921/kubernetes-client-go-creating-services-and-enpdoints
	service := apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "port-service",
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:     "service-port",
					Protocol: "TCP",
					Port:     8080,
				},
			},
		},
	}
	_, err := clientset.CoreV1().Services("default").Create(context.Background(), &service, metav1.CreateOptions{})
	if err != nil {
		log.Printf("cannot create service: %v", err)
	} else {
		fmt.Println("service created successfully.")
	}
}

func ListServices(serviceClient v1.ServiceInterface, clientset *kubernetes.Clientset) {
	list, err := serviceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list services: %v", err)
	}
	for _, d := range list.Items {
		fmt.Println(d.Name)
	}
}

func EditService(serviceClient v1.ServiceInterface, clientset *kubernetes.Clientset) {
	serviceName := "port-service"
	//reference: https://pkg.go.dev/k8s.io/client-go/util/retry#RetryOnConflict
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := serviceClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get latest version of Service: %v", err)
		}
		//fmt.Println(result)
		result.Spec.Ports[0].Port = 8090 // update port to 8090
		_, updateErr := serviceClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Printf("could not update service: %v", retryErr)
	} else {
		fmt.Println("service updated.")
	}
}

func DeleteService(serviceClient v1.ServiceInterface, clientset *kubernetes.Clientset) {
	serviceName := "port-service"
	deletePolicy := metav1.DeletePropagationForeground
	err := serviceClient.Delete(context.TODO(), serviceName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		log.Printf("could not delete service: %v", err)
	} else {
		fmt.Println("Service deleted.")
	}
}
