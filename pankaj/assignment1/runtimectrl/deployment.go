package runtimectrl

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pankaj9310/k8s-dev-training/pankaj/assignment1/goclient"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var deplymentName = "demo-deployment"

//DeploymentOperations - perform CURD operations on deployment object using runtime controller library
func DeploymentOperations() {

	ctrClient := getClient()

	fileBytes, err := ioutil.ReadFile("configfile/deployment.yaml")
	if err != nil {
		panic(err.Error())
	}

	var deploymentSpec appsv1.Deployment
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1024)
	err = dec.Decode(&deploymentSpec)
	if err != nil {
		panic(err.Error())
	}

	// Create Deployment

	fmt.Println("Deployment Create operation started ..........")
	err = ctrClient.Create(context.Background(), &deploymentSpec)
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(15 * time.Second) //wait to create deployment
	fmt.Println("Deployment Create operation completed ........")

	// Get Deployment
	fmt.Println("Deployment Get operation started ........")
	err = ctrClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: deplymentName}, &deploymentSpec)

	if err != nil {

		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err.Error()))
	}
	fmt.Printf("Latest Deployment is: %s with %d replicas \n", deploymentSpec.Name, *deploymentSpec.Spec.Replicas)
	fmt.Println("Deployment Get operation completed ........")

	// Update Deplyment
	fmt.Println("Deployment Update operation started .........")

	deploymentSpec.Spec.Replicas = goclient.ConvertToPtr(3)              // Increase replicaset count
	deploymentSpec.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" //change nginx version

	err = ctrClient.Update(context.Background(), &deploymentSpec)

	if err != nil {

		panic(fmt.Errorf("Failed to update Deployment: %v", err.Error()))
	}

	fmt.Println("Deployment update operation completed ..........")

	//Wait to update deplyment
	time.Sleep(30 * time.Second)

	// Verify Deployment
	fmt.Println("Deployment Update verfication operation started ........")
	err = ctrClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: deplymentName}, &deploymentSpec)

	if err != nil {

		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", err.Error()))
	}

	if *deploymentSpec.Spec.Replicas == 3 && deploymentSpec.Spec.Template.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Deployment update verfication successfull")

	} else {

		panic(fmt.Errorf("Deplyment verification failed. Replicas found: %d, expected 3 and Image found %s, expected: nginx:1.13",
			*deploymentSpec.Spec.Replicas, deploymentSpec.Spec.Template.Spec.Containers[0].Image))
	}
	fmt.Println("Deployment Update Verfication operation completed ........")

	// Delete Deployments
	fmt.Println("Deployment Delete operation started ........")
	err = ctrClient.Delete(context.TODO(), &deploymentSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Deployment Delete operation completed ........")
}
