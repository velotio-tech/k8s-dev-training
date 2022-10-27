package controllerruntime

import (
	"log"

	"github.com/velotio-tech/k8s-dev-training/services"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ControllerRuntime struct {
	client.Client
}

func GetClient() *ControllerRuntime {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Panicln(err)
		return nil
	}
	c, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Panicln(err)
		return nil
	}
	return &ControllerRuntime{c}
}

func (c *ControllerRuntime) GetResource(name string) services.Resource {
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
