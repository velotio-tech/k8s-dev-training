package service

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/clients"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/client_go/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GoClientTester() {
	var err error
	client := clients.NewServiceK8sGo(metav1.NamespaceDefault)
	listOptions := metav1.ListOptions{}
	serviceName := "nginx-service"

	// Create
	err = service.Create(client, "nginx-service", map[string]string{"hello": "world"}, 30080, 80, 80)
	if err != nil {
		fmt.Println("Failed to create service cuz", err)
	}

	// Get
	err = service.Get(client, listOptions)
	if err != nil {
		fmt.Println("Failed to get service cuz", err)
	}

	// Update
	err = service.Update(client, serviceName, map[string]string{"hello": "world"})
	if err != nil {
		fmt.Println("Failed to update service cuz", err)
	}

	// Delete
	err = service.Delete(client, serviceName)
	if err != nil {
		fmt.Println("Failed to delete service cuz", err)
	}
}
