package controllerruntimecrud

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateService(controllerClient client.Client) {
	// copy same body from client-go. This remains same. It won't change
	newService := apiv1.Service{
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
	err := controllerClient.Create(context.Background(), &newService)
	if err != nil {
		log.Printf("cannot create service: %v", err)
	} else {
		fmt.Println("Service created.")
	}

}
func ListServices(controllerClient client.Client) {
	services := &apiv1.ServiceList{}
	err := controllerClient.List(context.Background(), services)
	if err != nil {
		log.Printf("cannot access services list: %v", err)
	} else {
		for _, s := range services.Items {
			fmt.Println(s.Name)
		}
	}
}
func EditService(controllerClient client.Client) {
	service := &apiv1.Service{}
	err := controllerClient.Get(context.Background(), client.ObjectKey{
		Name: "port-service",
	}, service)
	if err != nil {
		log.Printf("cannot find desired service: %v", err)
	}
	service.Spec.Ports[0].Port = 8090 // // update port to 8090
	err = controllerClient.Update(context.TODO(), service)
	if err != nil {
		log.Printf("cannot update desired service: %v", err)
	} else {
		fmt.Println("Service updated.")
	}
}
func DeleteService(controllerClient client.Client) {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "port-service",
		},
	}
	err := controllerClient.Delete(context.Background(), service)
	if err != nil {
		log.Printf("cannot delete desired service: %v", err)
	} else {
		fmt.Println("Service deleted.")
	}
}
