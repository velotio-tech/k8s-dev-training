package clientgo

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
)

func Resources(clientset *kubernetes.Clientset) {
	ctx := context.Background()

	// SERVICES CRUD
	CreateService(clientset, ctx)
	ListServices(clientset, ctx)

	UpdateService(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListServices(clientset, ctx)

	DeleteService(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListServices(clientset, ctx)

	// DEPLOYMENTS CRUD
	CreateDeployment(clientset, ctx)

	ListDeployments(clientset, ctx)

	UpdateDeployment(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListDeployments(clientset, ctx)

	DeleteDeployment(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListDeployments(clientset, ctx)

	// CONFIGMAPS CRUD
	CreateConfigmap(clientset, ctx)
	ListConfigmaps(clientset, ctx)

	UpdateConfigmap(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListConfigmaps(clientset, ctx)

	DeleteConfigmap(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListConfigmaps(clientset, ctx)
}
