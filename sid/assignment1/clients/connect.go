package clients

import (
	"k8s.io/client-go/kubernetes"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func NewK8sGo() *kubernetes.Clientset {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientset
}

func NewPodsK8sGo(namespace string) clientcorev1.PodInterface {
	return NewK8sGo().CoreV1().Pods(namespace)
}

func NewDeployK8sGo(namespace string) clientappsv1.DeploymentInterface {
	return NewK8sGo().AppsV1().Deployments(namespace)
}
