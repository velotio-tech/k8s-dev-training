package deployment

import (
	"context"
	"fmt"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/common"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/utils/pointer"
	"time"
)

func Get(client appv1.DeploymentInterface, name string, options metav1.GetOptions) error {
	deployment, err := client.Get(context.TODO(), name, options)
	if err != nil {
		return err
	}
	fmt.Println(deployment.Name, deployment.Spec.Template.Spec.Containers[0].Image)
	return nil
}

func Create(client appv1.DeploymentInterface, name, image string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
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
							Name:  name,
							Image: image,
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
	_, err := client.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func Update(client appv1.DeploymentInterface, options metav1.UpdateOptions, name, image string) error {
	deployment, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	deployment.Spec.Template.Spec.Containers[0].Image = image
	_, err = client.Update(context.TODO(), deployment, options)
	if err != nil {
		return err
	}
	return nil
}

func Delete(client appv1.DeploymentInterface, name string, options metav1.DeleteOptions) error {
	return client.Delete(context.TODO(), name, options)
}

func Operations() error {
	ns := "default"
	client := common.DeploymentClient(ns)
	err := Create(client, "testdeploy", "nginx")
	if err != nil {
		return err
	}
	time.Sleep(15 * time.Second)
	err = Get(client, "testdeploy", metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = Update(client, metav1.UpdateOptions{}, "testdeploy", "nginx:1.17" )
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testdeploy", metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = Delete(client, "testdeploy", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testdeploy", metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}