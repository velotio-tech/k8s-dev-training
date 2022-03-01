package main

import (
	"assignment1/deployments"
	"assignment1/pods"
	"assignment1/services"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crtconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	controllerClient, err := client.New(crtconfig.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("Failed to create rtc client", err)
	}
	rtc := client.NewNamespacedClient(controllerClient, apiv1.NamespaceDefault)

	fmt.Println("Running using the Client GO library to create DEPLOYMENT, SERVICES, PODS")
	deployments.CreateDeploymentClient(clientset)
	services.CreateServicesClient(clientset)
	pods.SetPodsClient(clientset)

	fmt.Println("Creating Resources")
	deployments.CreateDeployment()
	services.CreateServices()
	pods.CreatePods()

	time.Sleep(20 * time.Second)
	fmt.Println("Created Following Resources")
	deployments.ListAllDeployments()
	services.GetAllServices()
	pods.ListAllPods()

	fmt.Println("Updating Resources")
	deployments.UpdateDeployment()
	services.UpdateServices()
	pods.UpdatePods()

	time.Sleep(20 * time.Second)
	fmt.Println("Deleting Resources")
	deployments.DeleteDeployment()
	services.DeleteService()
	pods.DeletePods()

	time.Sleep(40 * time.Second)

	fmt.Println("Running using the controller runtime library to create DEPLOYMENT, SERVICES, PODS")
	deployments.SetRtcClient(rtc)
	services.SetRtcClient(rtc)
	pods.SetRtcClient(rtc)

	fmt.Println("Creating Resources")
	deployments.CreateRtcDeployment()
	services.CreateRtcServices()
	pods.CreateRtcPods()

	time.Sleep(20 * time.Second)
	fmt.Println("Created Following Resources")
	deployments.ListRtcDeployments()
	services.ListRtcServices()
	pods.ListRtcPods()

	fmt.Println("Updating Resources")
	deployments.UpdateRtcDeployment()
	services.UpdateRtcService()
	pods.UpdateRtcPods()

	time.Sleep(20 * time.Second)
	fmt.Println("Deleting Resources")
	deployments.DeleteRtcDeployment()
	services.DeleteRtcService()
	pods.DeleteRtcPods()
}
