package deployment

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func Create(
	client clientappsv1.DeploymentInterface,
	replicas int32,
	labels map[string]string,
	name, image string,
) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: name + "-container", Image: image}},
				},
			},
		},
	}

	_, err := client.Create(context.TODO(), deployment, metav1.CreateOptions{})

	return err
}

func Get(client clientappsv1.DeploymentInterface, options metav1.ListOptions) error {
	deploymentList, err := client.List(context.TODO(), options)

	if err != nil {
		return err
	}

	for _, deployment := range deploymentList.Items {
		fmt.Println(deployment.Name)
	}

	return nil
}

func Update(client clientappsv1.DeploymentInterface, options metav1.ListOptions, labels map[string]string) error {
	deploymentList, err := client.List(context.TODO(), options)
	if err != nil {
		return err
	}

	for _, deployment := range deploymentList.Items {
		deployment.Labels = labels
		_, err = client.Update(context.TODO(), &deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func Delete(client clientappsv1.DeploymentInterface, options metav1.ListOptions) error {
	return client.DeleteCollection(context.TODO(), *metav1.NewDeleteOptions(int64(0)), options)
}
