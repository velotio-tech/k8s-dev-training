package main

import (
	"assignment1/deployments"
	"assignment1/pods"
	"assignment1/services"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crtconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Creating Deployment")
	deployments.CreateDeploymentClient(clientset)
	deployments.CreateDeployment()
	fmt.Println("Creating Services")
	services.CreateServicesClient(clientset)
	services.CreateServices()
	fmt.Println("Creating Pods")
	pods.CreatePodsClient(clientset)
	pods.CreatePods()

	time.Sleep(30 * time.Second)
	fmt.Println("Listing Deployment")
	deployments.ListDeployment()
	fmt.Println("Listing  Services")
	services.ListServices()
	fmt.Println("Listing Pods")
	pods.ListPods()

	time.Sleep(15 * time.Second)
	fmt.Println("Updating Deployments")
	deployments.UpdateDeployment()
	fmt.Println("Updating Services")
	services.UpdateServices()
	fmt.Println("Updating Pods")
	pods.UpdatePods()

	time.Sleep(30 * time.Second)
	fmt.Println("Listing Deployment")
	deployments.ListDeployment()
	fmt.Println("Listing  Services")
	services.ListServices()
	fmt.Println("Listing Pods")
	pods.ListPods()

	time.Sleep(15 * time.Second)
	fmt.Println("Delete Deployments")
	deployments.DeleteDeployment()
	fmt.Println("Delete Services")
	services.DeleteServices()
	fmt.Println("Delete Pods")
	pods.DeletePods()

	// creates the controller client
	controllerClient, err := client.New(crtconfig.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("Failed to create rtc client", err)
	}
	rtc := client.NewNamespacedClient(controllerClient, corev1.NamespaceDefault)

	fmt.Println("Creating controller runtime client for deployments, services and pods")
	deployments.CreateRTDeploymentClient(rtc)
	services.CreateRTServiceClient(rtc)
	pods.CreateRTPodClient(rtc)

	fmt.Println("Creating Deployment")
	deployments.CreateRTDeployment()
	fmt.Println("Creating Services")
	services.CreateRTServices()
	fmt.Println("Creating Pods")
	pods.CreateRTPods()

	time.Sleep(30 * time.Second)
	fmt.Println("Listing Deployment")
	deployments.ListRTDeployment()
	fmt.Println("Listing  Services")
	services.ListRTServices()
	fmt.Println("Listing Pods")
	pods.ListRTPods()

	time.Sleep(15 * time.Second)
	fmt.Println("Updating Deployments")
	deployments.UpdateRTDeployment()
	fmt.Println("Updating Services")
	services.UpdateServices()
	fmt.Println("Updating Pods")
	pods.UpdatePods()

	time.Sleep(30 * time.Second)
	fmt.Println("Listing Deployment")
	deployments.ListRTDeployment()
	fmt.Println("Listing  Services")
	services.ListRTServices()
	fmt.Println("Listing Pods")
	pods.ListRTPods()

	time.Sleep(15 * time.Second)
	fmt.Println("Delete Deployments")
	deployments.DeleteRTDeployment()
	fmt.Println("Delete Services")
	services.DeleteRTServices()
	fmt.Println("Delete Pods")
	pods.DeleteRTPods()

}
