package main

import (
	"fmt"
	"github.com/farkaskid/k8s-dev-training/assignment1/client"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers/client_go/pod"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	var err error
	clientset := client.NewK8sGo()
	v1Client := clientset.CoreV1()

	listOptions := metav1.ListOptions{}

	// Client Go
	// Create
	err = pod.CreatePod(v1Client, "nginx-pod", "nginx")
	if err != nil {
		fmt.Println("Failed to create pod cuz", err)
	}

	// Get
	err = pod.GetPods(v1Client, listOptions)
	if err != nil {
		fmt.Println("Error getting pods", err)
	}

	// Update Multiple
	err = pod.UpdatePods(v1Client, listOptions, map[string]string{"Hello": "World"})
	if err != nil {
		fmt.Println("Failed to update pods cuz", err)
	}

	// Delete Multiple
	err = pod.DeletePods(v1Client, listOptions)
	if err != nil {
		fmt.Println("Error getting pods", err)
	}
}
