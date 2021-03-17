package runtimectrl

import (
	"fmt"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func getClient() client.Client {
	ctrClient, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("Failed to create client")
		os.Exit(1)
	}
	return ctrClient
}
