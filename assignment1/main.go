package main

import (
	"fmt"
	"github.com/hatred09/k8s-dev-training/assignment1/clientgo"
	ctrl "github.com/hatred09/k8s-dev-training/assignment1/controller"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crtconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	controllerClient, err := client.New(crtconfig.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}
	rtc := client.NewNamespacedClient(controllerClient, apiv1.NamespaceDefault)
	deployClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	serviceClient := clientSet.CoreV1().Services(apiv1.NamespaceDefault)
	podClient := clientSet.CoreV1().Pods(apiv1.NamespaceDefault)

	fmt.Println("Running using the Client GO library to create DEPLOYMENT, SERVICES, PODS")

	fmt.Println("Creating Resources")
	err = clientgo.CreateDeployment(deployClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.CreateServices(serviceClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.CreatePods(podClient)
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("Created Following Resources")
	err = clientgo.ListAllDeployments(deployClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.GetAllServices(serviceClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.ListAllPods(podClient)
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("Updating Resources")
	err = clientgo.UpdateDeployment(deployClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.UpdateServices(serviceClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.UpdatePods(podClient)
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("Deleting Resources")
	err = clientgo.DeleteDeployment(deployClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.DeleteService(serviceClient)
	if err != nil {
		panic(err)
	}
	err = clientgo.DeletePods(podClient)
	if err != nil {
		panic(err)
	}

	time.Sleep(40 * time.Second)

	fmt.Println("Running using the controller runtime library to create DEPLOYMENT, SERVICES, PODS")

	fmt.Println("Creating Resources")
	err = ctrl.CreateRtcDeployment(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.CreateRtcServices(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.CreateRtcPods(rtc)
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("Created Following Resources")
	err = ctrl.ListRtcDeployments(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.ListRtcServices(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.ListRtcPods(rtc)
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("Updating Resources")
	err = ctrl.UpdateRtcDeployment(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.UpdateRtcService(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.UpdateRtcPods(rtc)
	if err != nil {
		panic(err)
	}

	time.Sleep(20 * time.Second)

	fmt.Println("Deleting Resources")
	err = ctrl.DeleteRtcDeployment(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.DeleteRtcService(rtc)
	if err != nil {
		panic(err)
	}
	err = ctrl.DeleteRtcPods(rtc)
	if err != nil {
		panic(err)
	}
}
