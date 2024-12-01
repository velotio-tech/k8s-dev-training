package clientgocrud

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)
	//Create deployment client for default namespace
var deploymentsClient v1.DeploymentInterface = clientset.AppsV1().Deployments("default")

func CreateDeployment() error {
	//Deploying a new deployment to cluster
	var replicas int32 = 2
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox-deployment",
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
							Name:    "busybox-deployment",
							Image:   "busybox:latest",
							Command: []string{"sleep", "100000"},
						},
					},
				},
			},
		},
	}

	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return nil
}
func ListDeployments() {
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list Deployment: %v", err)
	}
	for _, d := range list.Items {
		fmt.Println(d.Name)
	}
}
func EditDeployment() {
	deploymentName := "busybox-deployment"
	var updatedReplicaCount int32 = 1

	//reference: https://pkg.go.dev/k8s.io/client-go/util/retry#RetryOnConflict
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get latest version of Deployment: %v", err)
		}
		result.Spec.Replicas = &updatedReplicaCount
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Printf("could not update deployment: %v", retryErr)
	} else {
		fmt.Println("Deployment updated.")
	}
}
func DeleteDeployment() {
	deploymentName := "busybox-deployment"
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		log.Printf("could not delete deployment: %v", err)
	} else {
		fmt.Println("Deployment deleted.")
	}
}
