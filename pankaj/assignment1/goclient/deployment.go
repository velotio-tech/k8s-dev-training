package goclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var deplymentName = "demo-deployment"

//DeploymentOperations - perform CURD operations on deployment object using client-go library
func DeploymentOperations() {
	kubeclient := KubeConfig()
	deploymentsClient := kubeclient.AppsV1().Deployments(namespace)

	fileBytes, err := ioutil.ReadFile("configfile/deployment.yaml")
	if err != nil {
		panic(fmt.Errorf("Failed to read deployment yaml file: %v", err.Error()))
	}

	var deploymentSpec appsv1.Deployment
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1024)
	err = dec.Decode(&deploymentSpec)
	if err != nil {
		panic(fmt.Errorf("Failed to decode pod yaml data: %v", err.Error()))
	}

	// Create Deployment
	fmt.Println("Deployment Create operation started ..........")
	output, err := deploymentsClient.Create(context.TODO(), &deploymentSpec, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("Failed to create Deployment using go-client: %v", err.Error()))
	}
	time.Sleep(15 * time.Second) //wait to create deployment
	fmt.Printf("Deployment created %q \n", output.GetObjectMeta().GetName())
	fmt.Println("Deployment Create operation completed ..........")

	// Get Deployment
	fmt.Println("Deployment Get operation started ........")
	output, err = deploymentsClient.Get(context.TODO(), deplymentName, metav1.GetOptions{})

	if err != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err.Error()))
	}
	fmt.Printf("Latest Deployment is: %s with %d replicas \n", output.Name, *output.Spec.Replicas)
	fmt.Println("Deployment Get operation completed ..........")

	// Update Deplyment
	fmt.Println("Deployment Update operation started .........")

	output.Spec.Replicas = ConvertToPtr(3)                       // Increase replicaset count
	output.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" //change nginx version

	_, err = deploymentsClient.Update(context.TODO(), output, metav1.UpdateOptions{})

	if err != nil {
		panic(fmt.Errorf("Failed to update Deployment: %v", err.Error()))
	}

	fmt.Println("Deployment update operation completed ..........")

	//Wait to update deplyment
	time.Sleep(30 * time.Second)

	// Verify Deployment
	fmt.Println("Deployment Update verfication operation completed ..........")
	output, err = deploymentsClient.Get(context.TODO(), deplymentName, metav1.GetOptions{})

	if err != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err.Error()))
	}

	if *output.Spec.Replicas == 3 && output.Spec.Template.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Deployment update verfication successfull")

	} else {
		panic(fmt.Errorf("Deplyment verification failed. Replicas found: %d, expected 3 and Image found %s, expected: nginx:1.13",
			*output.Spec.Replicas, output.Spec.Template.Spec.Containers[0].Image))
	}
	fmt.Println("Deployment UPdate verfication  operation completed ..........")

	// Delete Deployments
	fmt.Println("Deployment Delete operation started ........")
	deletePolicy := metav1.DeletePropagationForeground
	err = deploymentsClient.Delete(context.TODO(), deplymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(fmt.Errorf("Failed to delete Deployment: %v", err.Error()))
	}
	fmt.Println("Deployment Delete operation completed ........")
}
