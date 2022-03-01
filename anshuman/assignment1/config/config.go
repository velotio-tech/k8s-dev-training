package config

import (
	"log"

	"k8s.io/client-go/kubernetes"
	av1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var apiObj v1.CoreV1Interface
var appApiObj av1.AppsV1Interface
var cl client.Client

func InitKubeConfig() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	//kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	apiObj = clientset.CoreV1()
	appApiObj = clientset.AppsV1()
}

func InitClient() {
	var err error
	cl, err = client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}
}

func GetClient() client.Client {
	return cl
}

func GetAPIObj() v1.CoreV1Interface {
	return apiObj
}

func GetAppAPIObj() av1.AppsV1Interface {
	return appApiObj
}
