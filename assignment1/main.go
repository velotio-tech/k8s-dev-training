package main

import (
	"assignment1/networkpolicy"
	"assignment1/pods"
	"assignment1/service"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	networkingv1 "k8s.io/client-go/kubernetes/typed/networking/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	//client-go
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// Pod CRUD
	creteOptions := metav1.CreateOptions{}
	pod, err := clientset.CoreV1().Pods("default").Create(ctx, pods.GetPodObject(), creteOptions)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	getOptions := metav1.GetOptions{}
	pod, err = clientset.CoreV1().Pods("default").Get(ctx, "nginx", getOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image before update: ", pod.Spec.Containers[0].Image)
	pod.Spec.Containers[0].Image = "nginx:1.22"
	updateOptions := metav1.UpdateOptions{}
	pod, err = clientset.CoreV1().Pods("default").Update(ctx, pod, updateOptions)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)

	pod, err = clientset.CoreV1().Pods("default").Get(ctx, "nginx", getOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image after update: ", pod.Spec.Containers[0].Image)
	deleteOptions := metav1.DeleteOptions{}
	err = clientset.CoreV1().Pods("default").Delete(ctx, "nginx", deleteOptions)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Pod deleted: %v\n", pod.Name)
	}

	// Service CRUD
	svc, err := clientset.CoreV1().Services("default").Create(ctx, service.GetServiceObject(), creteOptions)
	if err != nil {
		panic(err)
	}

	svc, err = clientset.CoreV1().Services("default").Get(ctx, svc.Name, getOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ports before update: ", svc.Spec.Ports)
	httpsPort := service.GetServicePort("https", 443)
	svc.Spec.Ports = append(svc.Spec.Ports, httpsPort)

	svc, err = clientset.CoreV1().Services("default").Update(ctx, svc, updateOptions)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("Ports after update: ", svc.Spec.Ports)
	err = clientset.CoreV1().Services("default").Delete(ctx, svc.Name, deleteOptions)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Service deleted: %v\n", svc.Name)
	}

	// Network policy CRUD

	networkingClientSet, err := networkingv1.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	networkPolicy, err := networkingClientSet.NetworkPolicies("default").Create(ctx, networkpolicy.GetNetworkPolicyObject(), creteOptions)
	if err != nil {
		panic(err)
	}

	networkPolicy, err = networkingClientSet.NetworkPolicies("default").Get(ctx, networkPolicy.Name, getOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PodSelectors before update: ", networkPolicy.Spec.PodSelector)
	networkPolicy.Spec.PodSelector.MatchLabels["role"] = "frontend"

	networkPolicy, err = networkingClientSet.NetworkPolicies("default").Update(ctx, networkPolicy, updateOptions)
	if err != nil {
		log.Fatal(err)
	}

	networkPolicy, err = networkingClientSet.NetworkPolicies("default").Get(ctx, networkPolicy.Name, getOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PodSelectors after update: ", networkPolicy.Spec.PodSelector)
	err = networkingClientSet.NetworkPolicies("default").Delete(ctx, networkPolicy.Name, deleteOptions)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("NetworkPolicy deleted: %v\n", networkPolicy.Name)
	}

	//controller-runtime

	fmt.Println("Using client from controller-runtime package")

	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}

	// Pod CRUD
	newPod := pods.GetPodObject()
	err = k8sClient.Create(ctx, newPod)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: newPod.Name}, newPod)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image before update: ", newPod.Spec.Containers[0].Image)

	newPod.Spec.Containers[0].Image = "nginx:1.22"

	err = k8sClient.Update(ctx, newPod)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)

	err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: newPod.Name}, newPod)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Image after update: ", newPod.Spec.Containers[0].Image)

	err = k8sClient.Delete(ctx, newPod)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Pod deleted: %v\n", pod.Name)
	}

	// Service CRUD

	newSvc := service.GetServiceObject()

	err = k8sClient.Create(ctx, newSvc)
	if err != nil {
		panic(err)
	}

	err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: newSvc.Name}, newSvc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ports before update: ", newSvc.Spec.Ports)

	newSvc.Spec.Ports = append(newSvc.Spec.Ports, httpsPort)

	err = k8sClient.Update(ctx, newSvc)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("Ports after update: ", newSvc.Spec.Ports)

	err = k8sClient.Delete(ctx, newSvc)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Service deleted: %v\n", newSvc.Name)
	}

	// Network policy CRUD

	newNetworkPolicy := networkpolicy.GetNetworkPolicyObject()
	err = k8sClient.Create(ctx, newNetworkPolicy)
	if err != nil {
		panic(err)
	}

	err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: newNetworkPolicy.Name}, newNetworkPolicy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PodSelectors before update: ", newNetworkPolicy.Spec.PodSelector)
	newNetworkPolicy.Spec.PodSelector.MatchLabels["role"] = "frontend"

	err = k8sClient.Update(ctx, newNetworkPolicy)
	if err != nil {
		log.Fatal(err)
	}

	err = k8sClient.Get(ctx, types.NamespacedName{Namespace: "default", Name: newNetworkPolicy.Name}, newNetworkPolicy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PodSelectors after update: ", newNetworkPolicy.Spec.PodSelector)
	err = k8sClient.Delete(ctx, newNetworkPolicy)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("NetworkPolicy deleted: %v\n", newNetworkPolicy.Name)
	}
}
