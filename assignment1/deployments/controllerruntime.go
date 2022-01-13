package deployments

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var rtc client.Client
var deploy = &appsv1.Deployment{}

func CreateRTDeploymentClient(rtClient client.Client) {
	rtc = rtClient
}

func CreateRTDeployment() {
	// Create Deployment
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-deploy",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "rtc-deploy-container",
							Image: "nginx",
							Ports: []v1.ContainerPort{
								{
									Name:          "port1",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println("Creating Deployment...")
	err := rtc.Create(context.Background(), deploy)
	if err != nil {
		log.Println("Error occcured while creating deploymnent", err.Error())
	}
	fmt.Println(deploy.Name, deploy.Spec.Template.Spec.Containers[0].Image)
}

func ListRTDeployment() {
	fmt.Println("Listing Deployments")
	deployments := &appsv1.DeploymentList{}
	err := rtc.List(context.Background(), deployments)
	if err != nil {
		fmt.Println("Error occured while fetching deployment list", err.Error())
	} else {
		for _, each := range deployments.Items {
			fmt.Println("Name : ", each.Name, " Lables : ", each.Labels)
		}
	}

}

func UpdateRTDeployment() {
	//Update Deployment

	deployment := &appsv1.Deployment{}
	err := rtc.Get(context.TODO(), client.ObjectKey{
		Name: "rtc-deploy",
	}, deployment)
	if err != nil {
		fmt.Println(err)
	}
	deployment.Spec.Template.Spec.Containers[0].Image = "nginx:1.17"
	//deployment.Spec.Replicas = int32Ptr(1)
	err = rtc.Update(context.TODO(), deployment)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(deployment.Name, deployment.Spec.Template.Spec.Containers[0].Image, deployment.Spec.Replicas)

}

//func int32Ptr(i int32) *int32 { return &i }

func DeleteRTDeployment() {
	//Delete Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-deploy",
		},
	}
	err := rtc.Delete(context.Background(), deployment)
	if err != nil {
		fmt.Println(err)
	}
}
