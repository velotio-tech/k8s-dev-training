package deployment

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/clients"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/controller-runtime/deployment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CRTTester() {
	var err error
	client := clients.NewCRTClient(metav1.NamespaceDefault)

	// Create
	err = deployment.Create(client, int32(2), map[string]string{"app": "nginx"}, "nginx", "nginx")
	if err != nil {
		fmt.Println("Failed to create deployment cuz", err)
	}

	// Get
	err = deployment.Get(client)
	if err != nil {
		fmt.Println("Failed to get deployments cuz", err)
	}

	// Update
	err = deployment.Update(client, "nginx", map[string]string{"hello": "world"})
	if err != nil {
		fmt.Println("Failed to update deployment cuz", err)
	}

	// Delete
	err = deployment.Delete(client, "nginx")
	if err != nil {
		fmt.Println("Failed to delete deployment cuz", err)
	}
}
