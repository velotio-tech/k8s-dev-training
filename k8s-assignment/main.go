package main

import (
	"fmt"
	assignment "k8s-assignment/assignment1"
	"k8s-assignment/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"time"
)

func main() {
	kubeconfig := utils.GetKubeConfig()
	coreclient := utils.CreateClient(kubeconfig)
	namespace := metav1.NamespaceDefault

	// Create deployment client
	deploymentsClient := coreclient.AppsV1().Deployments(namespace)
	// Create Deployment
	deploymentName := assignment.CreateDeployment(deploymentsClient)
	// Get Deployment spec
	fmt.Println(assignment.GetDeployment(deploymentsClient, deploymentName))
	// Update Deployment
	assignment.UpdateDeployment(deploymentsClient, deploymentName)
	// Get Deployment after updating spec
	fmt.Println(assignment.GetDeployment(deploymentsClient, deploymentName))
	fmt.Println("Deployment will be deleted in the 60 secs")
	time.Sleep(60 * time.Second)
	// Finally delete deployment
	assignment.DeleteDeployment(deploymentsClient, deploymentName)
	// Done with Client-go

	// Using controller-runtime to interact with k8s apis
	controllerClient := utils.CreateControllerClient(kubeconfig, &namespace)
	assignment.ListServices(controllerClient)
	serviceName := assignment.CreateService(controllerClient)
	assignment.EditServiceByName(controllerClient, serviceName)
	assignment.ListServices(controllerClient)

	fmt.Println("Service will be deleted in the 60 secs")
	time.Sleep(60 * time.Second)
	assignment.DeleteService(controllerClient, serviceName)
	assignment.ListServices(controllerClient)
}
