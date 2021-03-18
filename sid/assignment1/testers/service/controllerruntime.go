package service

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/clients"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/controller-runtime/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CRTTester() {
	var err error
	client := clients.NewCRTClient(metav1.NamespaceDefault)

	// Create
	err = service.Create(client, "nginx-service", map[string]string{"app": "nginx"}, 30080, 80, 80)
	if err != nil {
		fmt.Println("Failed to create service cuz", err)
	}

	// Get
	err = service.Get(client)
	if err != nil {
		fmt.Println("Failed to get services cuz", err)
	}

	// Update
	err = service.Update(client, "nginx-service", map[string]string{"hello": "world"})
	if err != nil {
		fmt.Println("Failed to update service cuz", err)
	}

	// Delete
	err = service.Delete(client, "nginx-service")
	if err != nil {
		fmt.Println("Failed to delete service cuz", err)
	}
}
