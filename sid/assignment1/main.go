package main

import (
	"github.com/farkaskid/k8s-dev-training/assignment1/client"
	"github.com/farkaskid/k8s-dev-training/assignment1/helpers"
)

func main() {
	clientset := client.New()
	helpers.GetPods(clientset)
}
