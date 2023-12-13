package clientgo

import (
	"assig1/constants"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateDeployment ...
func (cg *ClientGoClient) CreateDeployment() (*appsv1.Deployment, error) {
	replicas := int32(2)
	newDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cg.setUpEssentials.DeploymentName,
			Namespace: cg.setUpEssentials.Namespace,
			Labels: map[string]string{
				"app": "ariarijitAssign1",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "arijitAssign1",
				},
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "arijitAssign1",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: constants.NginxImage,
						},
					},
				},
			},
		},
	}

	createdDeployment, err := cg.deploymentClient.Create(context.Background(), newDeployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createdDeployment, nil
}

// UpdateDeployment ...
func (cg *ClientGoClient) UpdateDeployment(deploymentName string) error {
	updatedDeployment, err := cg.deploymentClient.Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting Pod for update: %v", err)
	}

	updatedDeployment.SetGenerateName("generated-updated-deployment")

	_, err = cg.deploymentClient.Update(context.Background(), updatedDeployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating Pod: %v", err)
	}

	return nil
}

// GetDeployments ...
func (cg *ClientGoClient) GetDeployments() (*appsv1.DeploymentList, error) {
	return cg.deploymentClient.List(context.Background(), metav1.ListOptions{})
}

// DeleteDeployment ...
func (cg *ClientGoClient) DeleteDeployment(deploymentName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	return cg.deploymentClient.Delete(context.Background(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}
