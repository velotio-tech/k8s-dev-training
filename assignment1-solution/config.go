package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func execute() {
	// Get the incluster config
	//config, err := rest.InClusterConfig() ------- TODO

	kubeconfig := flag.String("kubeconfig", "/home/priyankasalunke/.kube/config", "Location of kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		fmt.Println("Error occured during incluster config: ", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error occured during clientset: ", err)
	}

	resources(clientset)

}
