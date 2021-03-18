package clients

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func NewCRTClient(namespace string) client.Client {
	c, err := client.New(config.GetConfigOrDie(), client.Options{})

	if err != nil {
		fmt.Println("Failed to create CRT client cuz", err)
	}
	return client.NewNamespacedClient(c, namespace)
}
