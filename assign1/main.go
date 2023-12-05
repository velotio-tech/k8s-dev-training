package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	controllerConfig := ctrl.GetConfigOrDie()
	controllerClient, err := client.New(controllerConfig, client.Options{})
	if err != nil {
		panic(err.Error())
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "example-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}
	// Create a Pod
	createdPod, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, v1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Created Pod:", createdPod.Name)

	// Get a Pod
	foundPod, err := clientset.CoreV1().Pods("default").Get(context.Background(), "example-pod", v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Found Pod:", foundPod.Name)

	// Update a Pod
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		foundPod.Spec.Containers[0].Image = "nginx:latest"
		_, updateErr := clientset.CoreV1().Pods("default").Update(context.Background(), foundPod, v1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(retryErr.Error())
	}
	fmt.Println("Updated Pod image")

	// Delete a Pod
	deletePolicy := v1.DeletePropagationForeground
	if err := clientset.CoreV1().Pods("default").Delete(context.Background(), "example-pod", v1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err.Error())
	}
	fmt.Println("Deleted Pod")

	// Example operations using controller-runtime
	deployment := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      "example-deployment",
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(2),
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{"app": "example"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"app": "example"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
						},
					},
				},
			},
		},
	}

	// Create a Deployment using controller-runtime client
	if err := controllerClient.Create(context.Background(), deployment); err != nil {
		panic(err.Error())
	}
	fmt.Println("Created Deployment")

	// Get a Deployment using controller-runtime client
	foundDeployment := &appsv1.Deployment{}
	if err := controllerClient.Get(context.Background(), client.ObjectKey{Namespace: "default", Name: "example-deployment"}, foundDeployment); err != nil {
		panic(err.Error())
	}
	fmt.Println("Found Deployment:", foundDeployment.Name)

	// Update a Deployment using controller-runtime client
	foundDeployment.Spec.Replicas = func(i int32) *int32 { return &i }(3)
	if err := controllerClient.Update(context.Background(), foundDeployment); err != nil {
		panic(err.Error())
	}
	fmt.Println("Updated Deployment")

	// Delete a Deployment using controller-runtime client
	if err := controllerClient.Delete(context.Background(), foundDeployment); err != nil {
		panic(err.Error())
	}
	fmt.Println("Deleted Deployment")
}
