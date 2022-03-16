package config

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var clientSet *kubernetes.Clientset
var cl client.Client

func init() {
	// Init client
	var err error
	cl, err = client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}

	// Init Clientset
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	//kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// apiObj = clientSet.CoreV1()
	// appApiObj = clientSet.AppsV1()
}

func GetClient() client.Client {
	return cl
}

func GetClientSet() *kubernetes.Clientset {
	return clientSet
}
