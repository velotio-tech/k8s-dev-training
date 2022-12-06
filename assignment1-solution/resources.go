package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func resources(clientset *kubernetes.Clientset) {
	ctx := context.Background()

	// POD
	fmt.Println("***CREATE POD***")
	new_pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-pod-demo",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx-container-demo",
					Image: "nginx:latest",
				},
			},
		},
	}

	_, err := clientset.CoreV1().Pods("default").Create(ctx, new_pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error occured during pod creation: ", err)
		return
	}
	fmt.Println("Pod created successfully!")

	fmt.Println("***GET POD***")
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error : ", err)
	}

	fmt.Println("Total pods in the default namespace: ", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, "\t", pod.Status.Phase, "\t", pod.CreationTimestamp.Time)
	}

	fmt.Println("***UPDATE POD***")
	// add/edit new labels into the pod
	pod, err := clientset.CoreV1().Pods("default").Get(ctx, "nginx-pod-demo", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error while getting the pod: ", err)
		return
	}

	if pod.Labels == nil {
		// there are no labels on the pod
		pod.Labels = make(map[string]string)
	}

	// let's add/append new label
	pod.Labels["ass-1"] = "client-go-app-pod"

	_, err = clientset.CoreV1().Pods("default").Update(ctx, pod, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error during update pod: ", err)
		return
	}
	fmt.Println("Pod Label added successfully!")

	fmt.Println("***DELETE POD***")
	err = clientset.CoreV1().Pods("default").Delete(ctx, "nginx-pod-demo", *metav1.NewDeleteOptions(0))
	if err != nil {
		fmt.Println("Error occured in pod deletion")
		return
	}
	fmt.Println("Pod is deleted..!")

	// SERVICES

	fmt.Println("***CREATE SERVICE***")
	new_svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc-demo",
		},
		Spec: v1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	_, err = clientset.CoreV1().Services("default").Create(ctx, new_svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error occured during service creation, error: ", err)
		return
	}
	fmt.Println("Service created successfully!")

	fmt.Println("***GET SERVICE***")
	svc, err := clientset.CoreV1().Services("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error occured during get services: ", err)
		return
	}
	fmt.Println("Total services: ", len(svc.Items))

	for _, svc := range svc.Items {
		fmt.Println(svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, svc.Spec.Ports)
	}

	fmt.Println("***UPDATE SERVICE***")
	svc1, err := clientset.CoreV1().Services("default").Get(ctx, "svc-demo", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Service not found, error:", err)
		return
	}

	if svc1.Spec.Selector == nil {
		svc1.Spec.Selector = make(map[string]string)
	}

	svc1.Spec.Selector["app-1"] = "client-go-app-svc"
	_, err = clientset.CoreV1().Services("default").Update(ctx, svc1, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error occured during update service, error: ", err)
		return
	}
	fmt.Println("Service Updated successfully!")

	fmt.Println("***DELETE SERVICE***")
	err = clientset.CoreV1().Services("default").Delete(ctx, "svc-demo", metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Error occured during delete operation!", err)
		return
	}
	fmt.Println("Service deleted successfully!")

	// DEPLOYMENTS

	fmt.Println("***CREATE DEPLOYMENT***")
	var replicas int32 = 2
	new_deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment-demo",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "client-go-app",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "client-go-app",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:    "nginx-deployment",
							Image:   "nginx:latest",
							Command: []string{"sleep", "2000"},
						},
					},
				},
			},
		},
	}

	_, err = clientset.AppsV1().Deployments("default").Create(ctx, new_deploy, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error occured during deployment create: ", err)
		return
	}
	fmt.Println("Deployment created successfully")

	fmt.Println("***GET DEPLOYMENT***")
	deployment, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error occured in deployment list: ", err)
		return
	}

	fmt.Println("Total Deployments: ", len(deployment.Items))
	for _, deploy := range deployment.Items {
		fmt.Println(deploy.Name, deploy.Status.Replicas, deploy.Status.AvailableReplicas)
	}

	fmt.Println("***UPDATE DEPLOYMENT***")
	// get the deployment first and update the deployment labels
	dep, err := clientset.AppsV1().Deployments("default").Get(ctx, "nginx-deployment-demo", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Error occured during get deployment: ", err)
		return
	}

	if dep.Labels == nil {
		dep.Labels = make(map[string]string)
	}

	dep.Labels["ass-1"] = "client-go-app-deployment"
	_, err = clientset.AppsV1().Deployments("default").Update(ctx, dep, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error occured during update deployment: ", err)
		return
	}
	fmt.Println("Deployment updated successfully!")

	fmt.Println("***DELETE DEPLOYMENT***")
	err = clientset.AppsV1().Deployments("default").Delete(ctx, "nginx-deployment-demo", metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Error occured during deployment delete: ", err)
		return
	}
	fmt.Println("Deployment deleted successfully!")

}
