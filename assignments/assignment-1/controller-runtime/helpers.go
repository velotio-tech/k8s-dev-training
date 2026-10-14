package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getConfig() (*rest.Config, error) {
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if errors.Is(err, os.ErrNotExist) {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)
	}

	return config, nil
}
