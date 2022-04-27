package clientgocrud

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func CreatePod(podClient v1.PodInterface, clientset *kubernetes.Clientset) {

	fmt.Println("Deploying a pod to the cluster")
	newPod := &corev1.Pod{
		//this is metadata section of yaml file
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox-pod",
			Labels: map[string]string{
				"owner": "parav",
			},
		},
		Spec: corev1.PodSpec{
			//this is spec section of yaml file
			Containers: []corev1.Container{
				{Name: "busybox", Image: "busybox:latest", Command: []string{"sleep", "100000"}},
			},
		},
	}

	_, err := podClient.Create(context.Background(), newPod, metav1.CreateOptions{})
	if err != nil {
		log.Println("cannot create new pod: ", err)
	} else {
		fmt.Println("Pod deployed successfully.")
	}
}
func ListPods(podClient v1.PodInterface, clientset *kubernetes.Clientset) {
	fmt.Println("Listing Running Pods in the cluster")
	podList, err := podClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println("cannot get list of running pods:", err)
	}
	for _, n := range podList.Items {
		fmt.Println(n.Name)
	}
}
func EditPod(podClient v1.PodInterface, clientset *kubernetes.Clientset) {
	podname := "busybox-pod"
	updateOwner := "kaushal"
	result, err := podClient.Get(context.TODO(), podname, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get pod: %v", err)
	}
	if result.ObjectMeta.Labels["owner"] != "" {
		result.ObjectMeta.Labels["owner"] = updateOwner
	}

	_, err = podClient.Update(context.TODO(), result, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update pod: %v", err)
	} else {
		fmt.Println("Pod updated successfully.")
	}
}
func DeletePod(podClient v1.PodInterface, clientset *kubernetes.Clientset) {
	podname := "busybox-pod"
	err := podClient.Delete(context.TODO(), podname, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete pod: %v", err)
	} else {
		fmt.Println("Pod deleted successfully.")
	}
}
