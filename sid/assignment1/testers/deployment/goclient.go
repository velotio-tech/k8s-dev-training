package deployment

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/clients"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/client-go/deployment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GoClientTester() {
	var err error
	client := clients.NewDeployK8sGo(metav1.NamespaceDefault)
	listOptions := metav1.ListOptions{}

	// Create
	err = deployment.Create(client, int32(2), map[string]string{"app": "nginx"}, "nginx", "nginx")
	if err != nil {
		fmt.Println("Failed to create deployment cuz", err)
	}

	// Get
	err = deployment.Get(client, listOptions)
	if err != nil {
		fmt.Println("Failed to create deployment cuz", err)
	}

	// Update
	err = deployment.Update(client, listOptions, map[string]string{"Hello": "World"})
	if err != nil {
		fmt.Println("Failed to create deployment cuz", err)
	}

	// Delete
	err = deployment.Delete(client, listOptions)
	if err != nil {
		fmt.Println("Failed to create deployment cuz", err)
	}
}
