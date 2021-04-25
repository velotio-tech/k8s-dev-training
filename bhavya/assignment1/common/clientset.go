package common

import (
	"k8s.io/client-go/kubernetes"
	appv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func Getclientset() *kubernetes.Clientset {
	// creates the in-cluster config
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	panic(err.Error())
	//}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func DeploymentClient(ns string) appv1.DeploymentInterface {
	return Getclientset().AppsV1().Deployments(ns)
}

func PodClient(ns string) v1.PodInterface {
	return Getclientset().CoreV1().Pods(ns)
}

func StatefulsetClient(ns string) appv1.StatefulSetInterface {
	return Getclientset().AppsV1().StatefulSets(ns)
}

func ConfigMapClient(ns string) v1.ConfigMapInterface {
	return Getclientset().CoreV1().ConfigMaps(ns)
}
