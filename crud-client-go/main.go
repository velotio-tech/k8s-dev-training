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

var clientset *kubernetes.Clientset
var controllerClient client.Client
var Number int

// https://www.digitalocean.com/community/tutorials/understanding-init-in-go
func init() {

	// Creating client for client-go crud
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		log.Printf("Could not make client config: %v", err)
	} else {
		fmt.Println("client for client-go crud created")
	}
	clientset = kubernetes.NewForConfigOrDie(config)

	// Creating client for crud using controller-runtime
	controllerClient, err := client.New(controllerruntimeconfig.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Printf("could not create controller runtime client: %v", err)
	} else {
		fmt.Println("Controller runtime client created.")
	}
	//set namespace to default namespace
	controllerClient = client.NewNamespacedClient(controllerClient, apiv1.NamespaceDefault)
}
func main() {

	/* ######################## Cient-Go Section starts ######################## */

	//-------------------- POD SECTION ----------------------

	//passing client variable to other package only once
	clientgocrud.SetClient(clientset)
	// Creating a Pod
	if err := clientgocrud.CreatePod(); err != nil {
		fmt.Println("cannot create pod")
		panic(err)
	}
	// time.Sleep(10 * time.Second)
	//Listing all pods
	if err := clientgocrud.ListPods(); err != nil {
		fmt.Println("cannot list pods", err)
	}

	//Editing a Pod
	clientgocrud.EditPod()

	//Deleting a Pod
	clientgocrud.DeletePod()

	//-------------------- DEPLOYMENT SECTION ----------------------
	//Creating a deployment
	if err := clientgocrud.CreateDeployment(); err != nil {
		fmt.Println("cannot create deployment")
		panic(err)
	}

	//Listing all deployments
	clientgocrud.ListDeployments()

	//Editing a deployment
	clientgocrud.EditDeployment()

	//Deleting a deployment
	clientgocrud.DeleteDeployment()

	//-------------------- SERVICE SECTION ----------------------

	//Creating a service
	if err := clientgocrud.CreateService(); err != nil {
		fmt.Println("cannot create service")
		panic(err)
	}

	//Listing all services
	clientgocrud.ListServices()

	//Editing a service
	clientgocrud.EditService()

	//Deleting a service
	clientgocrud.DeleteService()

	/* ######################## Cient-Go Section ends ######################## */

	/* ######################## Controller-runtime Section starts ######################## */

	//-------------------- POD SECTION ----------------------
	//Creating a pod
	if err := controllerruntimecrud.CreatePod(controllerClient); err != nil {
		fmt.Println("cannot create pod")
		panic(err)
	}
	//Listing all pods
	controllerruntimecrud.ListPods(controllerClient)
	//Editing a pod
	controllerruntimecrud.EditPod(controllerClient)
	//Deleting a pod
	controllerruntimecrud.DeletePod(controllerClient)

	//-------------------- DEPLOYMENT SECTION ----------------------
	//Creating a deployment
	if err := controllerruntimecrud.CreateDeployment(controllerClient); err != nil {
		fmt.Println("cannot create deployment")
		panic(err)
	}
	//Listing all deployments
	controllerruntimecrud.ListDeployments(controllerClient)
	//Editing a deployment
	controllerruntimecrud.EditDeployment(controllerClient)
	//Deleting a deployment
	controllerruntimecrud.DeleteDeployment(controllerClient)

	//-------------------- SERVICE SECTION ----------------------
	//Creating a service
	if err := controllerruntimecrud.CreateService(controllerClient); err != nil {
		fmt.Println("cannot create service")
		panic(err)
	}
	//Listing all services
	controllerruntimecrud.ListServices(controllerClient)
	//Editing a service
	controllerruntimecrud.EditService(controllerClient)
	//Deleting a service
	controllerruntimecrud.DeleteService(controllerClient)

	/* ######################## CONTROLLER-RUNTIME SECTION ends ######################## */

}
