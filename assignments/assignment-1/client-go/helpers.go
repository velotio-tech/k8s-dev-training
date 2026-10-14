package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func createClientSet() (*kubernetes.Clientset, error) {

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

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}

	return clientSet, nil
}

func WaitForCondition[T runtime.Object](
	ctx context.Context,
	watcher watch.Interface,
	condition func(T) bool,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-watcher.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed")
			}

			if event.Type == watch.Error {
				return fmt.Errorf("watch error: %v", event.Object)
			}

			obj, ok := event.Object.(T)
			if !ok {
				continue
			}

			if condition(obj) {
				return nil
			}
		}
	}
}

func WaitForDeletion(
	ctx context.Context,
	watcher watch.Interface,
) error {
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-watcher.ResultChan():
			if !ok {
				return fmt.Errorf("watch channel closed")
			}

			if event.Type == watch.Error {
				return fmt.Errorf("watch error: %v", event.Object)
			}

			if event.Type == watch.Deleted {
				return nil
			}
		}
	}
}
