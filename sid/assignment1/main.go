package main

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/client"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/client_go/pod"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var err error
	podsGoClient := client.NewPodsK8sGo(metav1.NamespaceDefault)
	listOptions := metav1.ListOptions{}

	// Client Go
	// Create
	err = pod.CreatePod(podsGoClient, "nginx-pod", "nginx")
	if err != nil {
		fmt.Println("Failed to create pod cuz", err)
	}

	// Get
	err = pod.GetPods(podsGoClient, listOptions)
	if err != nil {
		fmt.Println("Error getting pods", err)
	}

	// Update
	err = pod.UpdatePods(podsGoClient, listOptions, map[string]string{"Hello": "World"})
	if err != nil {
		fmt.Println("Failed to update pods cuz", err)
	}

	// Delete
	err = pod.DeletePods(podsGoClient, listOptions)
	if err != nil {
		fmt.Println("Error getting pods", err)
	}
}
