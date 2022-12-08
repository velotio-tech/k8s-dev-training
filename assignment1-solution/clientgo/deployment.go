package clientgo

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateDeployment(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("***CREATE DEPLOYMENT***")
	var replicas int32 = 2
	new_deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment-demo",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "client-go-app",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "client-go-app",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "nginx-deployment",
							Image:   "nginx:latest",
							Command: []string{"sleep", "2000"},
						},
					},
				},
			},
		},
	}

	_, err := clientset.AppsV1().Deployments("default").Create(ctx, new_deploy, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error occured during deployment create: ", err)
		return
	}
	fmt.Println("Deployment created successfully")
	fmt.Println("--------------------------------")
}

func ListDeployments(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("***GET DEPLOYMENT***")
	deployment, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error occured in deployment list: ", err)
		return
	}

	fmt.Println("Total Deployments: ", len(deployment.Items))
	for _, deploy := range deployment.Items {
		fmt.Println(deploy.Name, deploy.Status.Replicas, deploy.Status.AvailableReplicas)
	}
	fmt.Println("--------------------------------")
}

func UpdateDeployment(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("***UPDATE DEPLOYMENT***")
	// get the deployment first and update the deployment labels
	dep, err := clientset.AppsV1().Deployments("default").Get(ctx, "proxy-app", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error occured during get deployment: ", err)
		return
	}

	if dep.Labels == nil {
		dep.Labels = make(map[string]string)
	}

	dep.Labels["ass-1"] = "client-go-app-deployment"
	_, err = clientset.AppsV1().Deployments("default").Update(ctx, dep, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error occured during update deployment: ", err)
		return
	}
	fmt.Println("Deployment updated successfully!")
	fmt.Println("--------------------------------")
}

func DeleteDeployment(clientset *kubernetes.Clientset, ctx context.Context) {

	fmt.Println("***DELETE DEPLOYMENT***")
	err := clientset.AppsV1().Deployments("default").Delete(ctx, "nginx-deployment-demo", metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Error occured during deployment delete: ", err)
		return
	}
	fmt.Println("Deployment deleted successfully!")
	fmt.Println("--------------------------------")
}
