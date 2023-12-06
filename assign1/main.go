package main

import (
	clientgo "assig1/client-go"
	crt "assig1/controller-runtime"
	"assig1/models"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
)

func main() {
	// client-go related ops
	go func() {
		specs := &models.ClientGoEssentials{
			Namespace:      corev1.NamespaceDefault,
			ServiceName:    "my-service",
			PodName:        "my-pod",
			DeploymentName: "my-deployment",
		}
		client, err := clientgo.NewClientGoClient(specs)
		if err != nil {
			log.Println("failed to initialize new client-go : ", err)
			return
		}

		// create pod
		createdPod, err := client.CreatePod()
		if err != nil {
			log.Println("failed to create the pod : ", err)
			return
		}
		log.Printf("created pod details:  %v\n", createdPod)

		// create service
		createService, err := client.CreateService()
		if err != nil {
			log.Println("failed to create the service : ", err)
			return
		}
		log.Printf("created service details:  %v\n", createService)

		// create deployment
		createdDeployment, err := client.CreateDeployment()
		if err != nil {
			log.Println("failed to create the deployment : ", err)
			return
		}
		log.Printf("created deploment details:  %v\n", createdDeployment)

		// get pods
		allPods, err := client.GetPods()
		if err != nil {
			log.Println("failed to get all pods : ", err)
			return
		}
		fmt.Printf("all pods : %v\n", allPods)

		allServices, err := client.GetServices()
		if err != nil {
			log.Println("failed to get all servcies : ", err)
			return
		}
		fmt.Printf("all services : %v\n", allServices)

		allDeployments, err := client.GetDeployments()
		if err != nil {
			log.Println("failed to get all deployments : ", err)
			return
		}
		fmt.Printf("all deployments : %v\n", allDeployments)

		// update ops
		err = client.UpdatePod(specs.PodName)
		if err != nil {
			log.Println("failed to update the pod : ", err)
			return
		}

		err = client.UpdateService(specs.ServiceName)
		if err != nil {
			log.Println("failed to update the servcie : ", err)
			return
		}

		err = client.UpdateDeployment(specs.DeploymentName)
		if err != nil {
			log.Println("failed to update the deployment : ", err)
			return
		}

		// delete ops
		err = client.DeletePod(specs.PodName)
		if err != nil {
			log.Println("failed to delete the pod : ", err)
			return
		}

		err = client.DeleteService(specs.ServiceName)
		if err != nil {
			log.Println("failed to delete the servcie : ", err)
			return
		}

		err = client.DeleteDeployment(specs.DeploymentName)
		if err != nil {
			log.Println("failed to delete the deployment : ", err)
			return
		}
	}()

	// controller-runtime related ops

	go func() {
		podName := "pod-crt"
		crtClient, err := crt.NewCRTClient(corev1.NamespaceDefault)
		if err != nil {
			log.Println("failed to initialize new client-go : ", err)
			return
		}

		err = crtClient.CreatePod(podName)
		if err != nil {
			log.Println("failed to create the pod : ", err)
			return
		}
		podDetails, err := crtClient.GetPod(podName)
		if err != nil {
			log.Println("failed to get the pod details : ", err)
			return
		}
		fmt.Printf("retrived pod details : %v\n", podDetails)

		err = crtClient.UpdatePod(podName)
		if err != nil {
			log.Println("failed to get the update pod details : ", err)
			return
		}
	}()
}
