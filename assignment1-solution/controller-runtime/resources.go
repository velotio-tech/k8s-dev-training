package controllerruntime

import (
	"context"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Resources(clientset client.Client) {
	ctx := context.Background()
	// set the default namespace otherwise it will throw an empty namespace error
	clientset = client.NewNamespacedClient(clientset, apiv1.NamespaceDefault)

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

	UpdateDeployments(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListDeployments(clientset, ctx)

	DeleteDeployments(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListDeployments(clientset, ctx)

	// CONFIGMAP CRUD
	CreateConfigmap(clientset, ctx)
	ListConfigmaps(clientset, ctx)

	UpdateConfigmap(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListConfigmaps(clientset, ctx)

	DeleteConfigmap(clientset, ctx)
	time.Sleep(5 * time.Second)
	ListConfigmaps(clientset, ctx)
}
