package util

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var clientset *kubernetes.Clientset = nil
var apiServerClient client.Client = nil

func GetInClusterKubeConfigClient() *kubernetes.Clientset {
	if clientset != nil {
		return clientset
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func GetCRTClient() (client.Client, error) {
	if apiServerClient == nil {
		var err error
		apiServerClient, err = client.New(config.GetConfigOrDie(), client.Options{})
		if err != nil {
			log.Println("error is not nil...")
			apiServerClient = nil
			return apiServerClient, err
		}
	}
	return apiServerClient, nil
}
