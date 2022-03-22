package clientgo

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

var deployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: "test-deployment",
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
						Image: "nginx",
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

func int32Ptr(i int32) *int32 { return &i }

func ListAllDeployments(deploymentClient v1.DeploymentInterface) error {

	fmt.Printf("Listing deployments in namespace %s:\n", apiv1.NamespaceDefault)
	list, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
	return nil
}

func CreateDeployment(deploymentClient v1.DeploymentInterface) error {

	fmt.Println("Creating deployment...", deploymentClient, deployment)
	result, err := deploymentClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %s.\n", result.GetObjectMeta().GetName())
	return nil
}

func DeleteDeployment(deploymentClient v1.DeploymentInterface) error {

	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentClient.Delete(context.Background(), "test-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	fmt.Println("Deleted deployment.")
	return nil
}

func UpdateDeployment(deploymentClient v1.DeploymentInterface) error {

	fmt.Println("Updating deployment...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentClient.Get(context.Background(), "test-deployment", metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		result.Spec.Replicas = int32Ptr(1)
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
		_, updateErr := deploymentClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}
	fmt.Println("Updated deployment...")
	return nil
}
