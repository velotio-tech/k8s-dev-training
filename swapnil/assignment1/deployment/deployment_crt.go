package main

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helpers "github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers"
	crcli "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	cli := helpers.GetCRTClient(metav1.NamespaceDefault)

	// creating pod
	fmt.Println("creating deployment")
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	err := cli.Create(context.TODO(), deployment)
	if err != nil {
		fmt.Println("Error creating the deployment", err)
	} else {
		// check if the pod is created
		deployment := &appsv1.Deployment{}
		err := cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-deployment"}, deployment)
		//err := cli.List(context.TODO(), podList)
		if err == nil {
			fmt.Println(deployment.Name, deployment.Labels)
		} else {
			fmt.Println("Error", err)
		}

		fmt.Println("config map creted")
	}
	fmt.Println("updating config map")
	deployment = &appsv1.Deployment{}
	err = cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-deployment"}, deployment)
	if err != nil {
		fmt.Println("Error fetching deployment", err)
	}else{
		fmt.Println("updating the deployment")
		//deployment.Labels["APPNAME"] = "Awesome"
		deployment.Spec.Replicas = int32Ptr(3)
		err := cli.Update(context.TODO(), deployment)
		if err != nil {
			fmt.Println("Error updating the deployment", err)
		} else {
			// verify if the lables are added
			deployment = &appsv1.Deployment{}
			err := cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-deployment"}, deployment)
			//err := cli.List(context.TODO(), podList)
			if err == nil {
				fmt.Println(deployment.Name, deployment.Labels)
			} else {
				fmt.Println("Error", err)
			}
		}
	}

	// listing deployments
	deploymentList := &appsv1.DeploymentList{}
	err = cli.List(context.TODO(), deploymentList)
	if err != nil {
		fmt.Println("error fetching deployment list")
	} else {
		for _, each := range deploymentList.Items {
			fmt.Println("name: ", each.Name, " Lables: ", each.Labels)
		}
	}

	// deleting the pod
	deployment = &appsv1.Deployment{}
	err = cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-deployment"}, deployment)
	if err != nil {
		fmt.Println("Error fetching deployment", err)
	} else {
		//cli.Delete(
		//	context.TODO(),
		//	deployment,
		//	&crcli.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))})

		fmt.Println("config mpa deleted")
	}

}

func int32Ptr(i int32) *int32 { return &i }