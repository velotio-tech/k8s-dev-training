package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateRtcDeployment(rtc client.Client) error {

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-deployment",
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
							Name:  "rtc-deployment-container",
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
		return err
	}
	fmt.Println(deploy.Name, deploy.Spec.Template.Spec.Containers[0].Image)
	return nil
}

func ListRtcDeployments(rtc client.Client) error {

	fmt.Println("Listing Deployments")
	deployments := &appsv1.DeploymentList{}
	err := rtc.List(context.Background(), deployments)
	if err != nil {
		return err
	} else {
		for _, each := range deployments.Items {
			fmt.Println("Name : ", each.Name, " Labels : ", each.Labels)
		}
	}
	return nil
}

func UpdateRtcDeployment(rtc client.Client) error {

	deployment := &appsv1.Deployment{}
	err := rtc.Get(context.TODO(), client.ObjectKey{
		Name: "rtc-deployment",
	}, deployment)
	if err != nil {
		return err
	}
	deployment.Spec.Template.Spec.Containers[0].Image = "nginx:1.17"
	err = rtc.Update(context.TODO(), deployment)
	if err != nil {
		return err
	}
	fmt.Println(deployment.Name, deployment.Spec.Template.Spec.Containers[0].Image, deployment.Spec.Replicas)
	return nil
}

func DeleteRtcDeployment(rtc client.Client) error {

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-deployment",
		},
	}
	err := rtc.Delete(context.Background(), deployment)
	if err != nil {
		return err
	}
	return nil
}
