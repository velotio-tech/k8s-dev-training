package main

import (
	"fmt"
	"log"

	controllerruntimecrud "github.com/paravkaushal/crud-client-go/controller-runtime-crud"
	apiv1 "k8s.io/api/core/v1"

	clientgocrud "github.com/paravkaushal/crud-client-go/client-go-crud"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	controllerruntimeconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {

	/* ######################## Cient-Go Section starts ######################## */

	// Creating client for client-go crud
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Printf("Could not make client config: %v", err)
	} else {
		fmt.Println("client for client-go crud created")
	}
	clientset := kubernetes.NewForConfigOrDie(config)

	//-------------------- POD SECTION ----------------------

	//Create pod client for default namespace
	podClient := clientset.CoreV1().Pods("default")
	// Creating a Pod
	clientgocrud.CreatePod(podClient, clientset)
	// time.Sleep(10 * time.Second)
	//Listing all pods
	clientgocrud.ListPods(podClient, clientset)

	//Editing a Pod
	clientgocrud.EditPod(podClient, clientset)

	//Deleting a Pod
	clientgocrud.DeletePod(podClient, clientset)

	//-------------------- DEPLOYMENT SECTION ----------------------
	//Create deployment client for default namespace
	deploymentsClient := clientset.AppsV1().Deployments("default")

	//Creating a deployment
	clientgocrud.CreateDeployment(deploymentsClient, clientset)

	//Listing all deployments
	clientgocrud.ListDeployments(deploymentsClient, clientset)

	//Editing a deployment
	clientgocrud.EditDeployment(deploymentsClient, clientset)

	//Deleting a deployment
	clientgocrud.DeleteDeployment(deploymentsClient, clientset)

	//-------------------- SERVICE SECTION ----------------------
	//Create service client for default namespace
	serviceClient := clientset.CoreV1().Services("default")

	//Creating a service
	clientgocrud.CreateService(serviceClient, clientset)

	//Listing all services
	clientgocrud.ListServices(serviceClient, clientset)

	//Editing a service
	clientgocrud.EditService(serviceClient, clientset)

	//Deleting a service
	clientgocrud.DeleteService(serviceClient, clientset)

	/* ######################## Cient-Go Section ends ######################## */

	/* ######################## Controller-runtime Section starts ######################## */

	fmt.Println("Going to create controller runtime client")
	controllerClient, err := client.New(controllerruntimeconfig.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Printf("could not create controller runtime client: %v", err)
	} else {
		fmt.Println("Controller runtime client created.")
	}
	//set namespace to default namespace
	controllerClient = client.NewNamespacedClient(controllerClient, apiv1.NamespaceDefault)

	//-------------------- POD SECTION ----------------------
	//Creating a pod
	controllerruntimecrud.CreatePod(controllerClient)
	//Listing all pods
	controllerruntimecrud.ListPods(controllerClient)
	//Editing a pod
	controllerruntimecrud.EditPod(controllerClient)
	//Deleting a pod
	controllerruntimecrud.DeletePod(controllerClient)

	//-------------------- DEPLOYMENT SECTION ----------------------
	//Creating a deployment
	controllerruntimecrud.CreateDeployment(controllerClient)
	//Listing all deployments
	controllerruntimecrud.ListDeployments(controllerClient)
	//Editing a deployment
	controllerruntimecrud.EditDeployment(controllerClient)
	//Deleting a deployment
	controllerruntimecrud.DeleteDeployment(controllerClient)

	//-------------------- SERVICE SECTION ----------------------
	//Creating a service
	controllerruntimecrud.CreateService(controllerClient)
	//Listing all services
	controllerruntimecrud.ListServices(controllerClient)
	//Editing a service
	controllerruntimecrud.EditService(controllerClient)
	//Deleting a service
	controllerruntimecrud.DeleteService(controllerClient)

	/* ######################## CONTROLLER-RUNTIME SECTION ends ######################## */

}
