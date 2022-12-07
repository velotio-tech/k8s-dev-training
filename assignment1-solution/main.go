package main

import (
	goclient "ass1/clientgo"
	controllerclient "ass1/controller-runtime"
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/priyankasalunke/.kube/config", "Location of kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Println("Error occured during incluster config: ", err)
	}

	fmt.Println("*** Client-Go ***")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error occured during clientset: ", err)
	}

	goclient.Resources(clientset)

	fmt.Println("*** Controller Runtime ***")

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}

	controllerclient.Resources(k8sClient)
}
