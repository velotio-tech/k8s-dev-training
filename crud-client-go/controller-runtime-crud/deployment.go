package controllerruntimecrud

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(controllerClient client.Client) {

	//same content from client-go file
	var replicas int32 = 2
	newDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox-deployment2",
			Labels: map[string]string{
				"owner": "parav",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"owner": "parav",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"owner": "parav",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "busybox-deployment2",
							Image:   "busybox:latest",
							Command: []string{"sleep", "100000"},
						},
					},
				},
			},
		},
	}

	err := controllerClient.Create(context.Background(), newDeployment)
	if err != nil {
		log.Printf("could not deploy Deployment: %v", err)
	} else {
		fmt.Println("Deployment Created")
	}

}
func ListDeployments(controllerClient client.Client) {
	deploymentList := &appsv1.DeploymentList{}
	err := controllerClient.List(context.Background(), deploymentList, client.InNamespace("default"))
	if err != nil {
		log.Printf("failed to list Deployments in namespace default: %v\n", err)
	} else {
		for _, d := range deploymentList.Items {
			fmt.Println(d.Name)
		}
	}
}
func EditDeployment(controllerClient client.Client) {
	deployment := &appsv1.Deployment{}
	err := controllerClient.Get(context.TODO(), client.ObjectKey{
		Name: "busybox-deployment2",
	}, deployment)
	if err != nil {
		log.Printf("could not find required deployment: %v", err)
	}
	deployment.ObjectMeta.Labels["owner"] = "kaushal"

	err = controllerClient.Update(context.TODO(), deployment)
	if err != nil {
		log.Printf("could not update required deployment: %v", err)
	} else {
		fmt.Println("Deployment Updated")
	}
}
func DeleteDeployment(controllerClient client.Client) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox-deployment2",
		},
	}
	err := controllerClient.Delete(context.Background(), deployment)
	if err != nil {
		log.Printf("could not delete required deployment: %v", err)
	} else {
		fmt.Println("Deployment deleted.")
	}
}
