package controllerruntime

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(clientset client.Client, ctx context.Context) {
	var replicas int32 = 3
	fmt.Println("*** CREATE DEPLOYMENT ***")
	new_deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment-demo",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "controller-runtime-app",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "controller-runtime-app",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Image: "nginx:latest",
							Name:  "nginx-controller-deploy",
						},
					},
				},
			},
		},
	}

	err := clientset.Create(ctx, new_deployment)
	if err != nil {
		fmt.Println("Deployment creation failed, error: ", err)
		return
	}
	fmt.Println("Deployment created successfully!")
	fmt.Println("--------------------------------")
}

func ListDeployments(clientset client.Client, ctx context.Context) {
	fmt.Println("*** GET DEPLOYMENT ***")
	deployment := &appsv1.DeploymentList{}
	err := clientset.List(ctx, deployment)
	if err != nil {
		fmt.Println("Error during listing deployments, error: ", err)
		return
	}
	fmt.Println("Total deployment in the default namespace: ", len(deployment.Items))

	for _, deploy := range deployment.Items {
		fmt.Println(deploy.Name, deploy.Status.Replicas, deploy.Status.AvailableReplicas)
	}
	fmt.Println("--------------------------------")
}

func UpdateDeployments(clientset client.Client, ctx context.Context) {
	fmt.Println("*** UPDATE DEPLOYMENT ***")
	deploy := &appsv1.Deployment{}
	err := clientset.Get(ctx, client.ObjectKey{Name: "nginx"}, deploy)
	if err != nil {
		fmt.Println("GET Deployment failed with an error: ", err)
		return
	}

	if deploy.ObjectMeta.Labels == nil {
		deploy.ObjectMeta.Labels = make(map[string]string)
	}

	deploy.ObjectMeta.Labels["ass-1"] = "deployment-label-demo"
	err = clientset.Update(ctx, deploy)
	if err != nil {
		fmt.Println("Update deployment is failed, error: ", err)
		return
	}
	fmt.Println("Deployment updated successfully!")
	fmt.Println("--------------------------------")
}

func DeleteDeployments(clientset client.Client, ctx context.Context) {
	fmt.Println("*** DELETE DEPLOYMENT ***")
	del_deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment-demo",
		},
	}

	err := clientset.Delete(ctx, del_deployment)
	if err != nil {
		fmt.Println("Failed to delete the deployment, error: ", err)
		return
	}
	fmt.Println("Deployment deleted successfully!")
	fmt.Println("--------------------------------")
}
