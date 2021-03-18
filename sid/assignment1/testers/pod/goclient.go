package pod

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/clients"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/client-go/pod"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GoClientTester() {
	var err error
	client := clients.NewPodsK8sGo(metav1.NamespaceDefault)
	listOptions := metav1.ListOptions{}

	// Client Go
	// Create
	err = pod.Create(client, "nginx-pod", "nginx")
	if err != nil {
		fmt.Println("Failed to create pod cuz", err)
	}

	// Get
	err = pod.Get(client, listOptions)
	if err != nil {
		fmt.Println("Error getting pods", err)
	}

	// Update
	err = pod.Update(client, listOptions, map[string]string{"Hello": "World"})
	if err != nil {
		fmt.Println("Failed to update pods cuz", err)
	}

	// Delete
	err = pod.Delete(client, listOptions)
	if err != nil {
		fmt.Println("Error getting pods", err)
	}
}
