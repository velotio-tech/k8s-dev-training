package utils

import (
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetKubeConfig() *string {
	// Get kubectl config
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return kubeconfig

}
func CreateClient(kubeconfig *string) *kubernetes.Clientset {
	// Create kubernets client set to interact with resources
	// This returns base client and derived clients can be created using this client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	coreclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return coreclient
}

func CreateControllerClient(kubeconfig *string, ns *string) client.Client {
	// Create kubernets client set to interact with resources
	// This returns base client and derived clients can be created using this client
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	controllerClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}
	controllerNsClient := client.NewNamespacedClient(controllerClient, *ns)
	return controllerNsClient
}
