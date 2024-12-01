package crt

import (
	"assig1/constants"
	"context"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreatePod ...
func (c *CRTClient) CreatePod(podName string) error {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: c.namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Containers: []corev1.Container{
				{
					Name:            "nginx-crt",
					Image:           constants.NginxImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
		},
	}
	return c.apiServerClient.Create(context.Background(), pod, &client.CreateOptions{})
}

// UpdatePod ...
func (c *CRTClient) UpdatePod(podName string) error {
	oldDeployment := &appsv1.Deployment{}
	err := c.apiServerClient.Get(context.Background(), client.ObjectKey{Namespace: c.namespace, Name: podName}, oldDeployment)
	if err != nil {
		return err
	}
	oldDeployment.SetGenerateName("updated-ctr-deployment")
	oldDeployment.Spec.Replicas = func(i int32) *int32 { return &i }(3)
	err = c.apiServerClient.Update(context.Background(), oldDeployment, &client.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// DeletePod ...
func (c *CRTClient) DeletePod(podName string) error {
	existingDeployment := &appsv1.Deployment{}
	err := c.apiServerClient.Get(context.Background(), client.ObjectKey{Namespace: c.namespace, Name: podName}, existingDeployment)
	if err != nil {
		return err
	}
	return c.apiServerClient.Delete(context.Background(), existingDeployment)
}

// GetPod ...
func (c *CRTClient) GetPod(podName string) (*appsv1.Deployment, error) {
	oldDeployment := &appsv1.Deployment{}
	err := c.apiServerClient.Get(context.Background(), client.ObjectKey{Namespace: c.namespace, Name: podName}, oldDeployment)
	if err != nil {
		return nil, err
	}

	return oldDeployment, nil
}
