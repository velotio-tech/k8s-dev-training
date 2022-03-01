package c_go

import (
	"assignment1/config"
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(name, namespace, image string, replicas int32) error {

	apiObj := config.GetAppAPIObj()

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

	createOptions := metav1.CreateOptions{}

	_, err := apiObj.Deployments(namespace).Create(context.Background(), &deployment, createOptions)

	return err
}

func DeleteDeployment(name, namespace string) error {

	apiObj := config.GetAppAPIObj()

	deleteOptions := metav1.DeleteOptions{}

	err := apiObj.Deployments(namespace).Delete(context.Background(), name, deleteOptions)

	return err
}

func ReadDeployment(name, namespace string) error {

	var showDeployments bool
	if name == "" {
		showDeployments = true
	}

	apiObj := config.GetAppAPIObj()

	if showDeployments {

		listOptions := metav1.ListOptions{}

		deployments, err := apiObj.Deployments(namespace).List(context.Background(), listOptions)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(deployments, "", "\t")
		if err != nil {
			return err
		}

		log.Println(string(b))
	} else {

		getOptions := metav1.GetOptions{}

		deployment, err := apiObj.Deployments(namespace).Get(context.Background(), name, getOptions)
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

	apiObj := config.GetAppAPIObj()

	getOptions := metav1.GetOptions{}

	deployment, err := apiObj.Deployments(namespace).Get(context.Background(), name, getOptions)
	if err != nil {
		return err
	}

	if image != "" {
		deployment.Spec.Template.Spec.Containers[0].Image = image
	}
	if replicas != 0 {
		deployment.Spec.Replicas = &replicas
	}

	updateOptions := metav1.UpdateOptions{}

	_, err = apiObj.Deployments(namespace).Update(context.Background(), deployment, updateOptions)

	return err

}
