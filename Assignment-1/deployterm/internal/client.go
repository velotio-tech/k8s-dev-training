package internal

import (
	"os"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sClient struct {
	clientset *kubernetes.Clientset
	runclient client.Client
}

func NewK8sClient(kubeconfigPath string) (*K8sClient, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ctrlC, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}

	return &K8sClient{clientset: clientset, runclient: ctrlC}, nil
}

func (client *K8sClient) GetClientset() *kubernetes.Clientset {
	return client.clientset
}

func (client *K8sClient) GetRunclient() client.Client {
	return client.runclient
}
