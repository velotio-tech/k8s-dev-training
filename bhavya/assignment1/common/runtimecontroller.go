package common

import (
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func RcontrollerClient(ns string) client.Client {
	clt, err := client.New(config.GetConfigOrDie(),client.Options{})
	if err != nil {
		fmt.Println("Failed to create rtc client", err)
	}
	return client.NewNamespacedClient(clt,ns)
}