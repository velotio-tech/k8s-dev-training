package goclient

import (
	"log"

	"github.com/velotio-tech/k8s-dev-training/services"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GoClient struct {
	*kubernetes.Clientset
}

func GetClient() *GoClient {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicln(err)
		return nil
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln(err)
		return nil
	}
	return &GoClient{clientSet}
}

func (c *GoClient) GetResource(name string) services.Resource {
	switch name {
	case "namespace", "ns", "namespaces":
		return &Namespace{c}
	case "pod", "pods", "po":
		return &Pod{c}
	case "svc", "service", "services":
		return &Service{c}
	}
	return nil
}
