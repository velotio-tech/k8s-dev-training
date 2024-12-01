package pod

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/clients"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/controller-runtime/pod"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CRTTester() {
	var err error
	client := clients.NewCRTClient(metav1.NamespaceDefault)

	// Get
	err = pod.Get(client)
	if err != nil {
		fmt.Println("Failed to get pods cuz", err)
	}

	// Create
	err = pod.Create(client, "nginx-pod", "nginx")
	if err != nil {
		fmt.Println("Failed to create pod cuz", err)
	}

	// Update
	err = pod.Update(client, "nginx-pod", map[string]string{"hello": "world"})
	if err != nil {
		fmt.Println("Failed to update pod cuz", err)
	}

	// Delete
	err = pod.Delete(client, "nginx-pod")
	if err != nil {
		fmt.Println("Failed to update pod cuz", err)
	}
}
