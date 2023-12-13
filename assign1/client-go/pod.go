package clientgo

import (
	"assig1/constants"
	"context"
	"fmt"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreatePod ...
func (cg *ClientGoClient) CreatePod() (*core.Pod, error) {
	newPod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cg.setUpEssentials.PodName,
			Namespace: cg.setUpEssentials.Namespace,
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "nginx",
					Image:           constants.NginxImage,
					ImagePullPolicy: core.PullIfNotPresent,
				},
			},
		},
	}

	createdPod, err := cg.podClient.Create(context.Background(), newPod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return createdPod, nil
}

// UpdatePod ...
func (cg *ClientGoClient) UpdatePod(podName string) error {
	updatedPod, err := cg.podClient.Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error getting Pod for update: %v", err)
	}

	updatedPod.SetGenerateName("generated-updated-name")
	updatedPod.Spec.Containers[0].Image = "nginx:1.21"

	_, err = cg.podClient.Update(context.Background(), updatedPod, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating Pod: %v", err)
	}

	return nil
}

// GetPods ...
func (cg *ClientGoClient) GetPods() (*core.PodList, error) {
	return cg.podClient.List(context.Background(), metav1.ListOptions{})
}

// DeletePod ...
func (cg *ClientGoClient) DeletePod(podName string) error {
	return cg.podClient.Delete(context.Background(), podName, metav1.DeleteOptions{})
}
