package deployment

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(client crtclient.Client, replicas int32, labels map[string]string, name, image string) error {
	return client.Create(context.TODO(), &appsv1.Deployment{
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
	})
}

func Get(client crtclient.Client) error {
	deploymentList := &appsv1.DeploymentList{}

	err := client.List(context.TODO(), deploymentList)
	if err != nil {
		return err
	}

	for _, deployment := range deploymentList.Items {
		fmt.Println(deployment.Name)
	}
	return nil
}

func Update(client crtclient.Client, name string, labels map[string]string) error {
	deployment := &appsv1.Deployment{}
	err := client.Get(context.TODO(), crtclient.ObjectKey{Name: name}, deployment)
	if err != nil {
		return err
	}

	deployment.Labels = labels

	return client.Update(context.TODO(), deployment)
}

func Delete(client crtclient.Client, name string) error {
	return client.Delete(
		context.TODO(),
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: name}},
		&crtclient.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))},
	)
}
