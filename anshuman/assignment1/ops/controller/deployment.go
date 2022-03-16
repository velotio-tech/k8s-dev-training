package controller

import (
	"assignment1/config"
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(name, namespace, image string, replicas int32) error {

	cl := config.GetClient()

	deployment := v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
			Namespace:    namespace,
			Labels:       map[string]string{"app": "deployment"},
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    map[string]string{"app": "deployment"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  name,
						Image: image,
					}},
				},
			},
		},
		Status: v1.DeploymentStatus{},
	}

	createOptions := client.CreateOptions{}

	err := cl.Create(context.Background(), &deployment, &createOptions)

	return err
}

func DeleteDeployment(name, namespace string) error {

	cl := config.GetClient()
	deployment := v1.Deployment{}
	objectKey := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	deleteOptions := client.DeleteOptions{}

	cl.Get(context.Background(), objectKey, &deployment)

	err := cl.Delete(context.Background(), &deployment, &deleteOptions)

	return err
}

func ReadDeployment(name, namespace string) error {

	var showDeployments bool
	if name == "" {
		showDeployments = true
	}

	cl := config.GetClient()

	if showDeployments {

		listOptions := client.ListOptions{
			Namespace: namespace,
		}
		deployments := v1.DeploymentList{}

		err := cl.List(context.Background(), &deployments, &listOptions)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(deployments, "", "\t")
		if err != nil {
			return err
		}

		log.Println(string(b))
	} else {

		deployment := v1.Deployment{}
		getOptions := client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}

		err := cl.Get(context.Background(), getOptions, &deployment)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(deployment, "", "\t")
		if err != nil {
			return err
		}

		log.Println(string(b))
	}
	return nil
}

func UpdateDeployment(name, namespace, image string, replicas int32) error {

	cl := config.GetClient()

	getOptions := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	updateOptions := client.UpdateOptions{}
	deployment := v1.Deployment{}

	err := cl.Get(context.Background(), getOptions, &deployment)
	if err != nil {
		return err
	}

	if image != "" {
		deployment.Spec.Template.Spec.Containers[0].Image = image
	}
	if replicas >= 0 {
		deployment.Spec.Replicas = &replicas
	}

	err = cl.Update(context.Background(), &deployment, &updateOptions)

	return err

}
