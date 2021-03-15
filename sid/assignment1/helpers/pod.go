package helpers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPods(clientset *kubernetes.Clientset) {
	coreInterface := clientset.CoreV1()
	pods, err := coreInterface.Pods("default").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	printPods(pods)
}

func printPods(pods *v1.PodList) {
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}
