package pods

import (
	"context"
	"fmt"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

var pods = &apiv1.Pod{
	TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
	ObjectMeta: metav1.ObjectMeta{Name: "demo-pod"},
	Spec:       apiv1.PodSpec{Containers: []apiv1.Container{{Name: "demo-nginx", Image: "nginx"}}},
}
var podsClient corev1.PodInterface

func CreatePodsClient(clientset *kubernetes.Clientset) {
	podsClient = clientset.CoreV1().Pods(apiv1.NamespaceDefault)
}
func CreatePods() {
	// Create Pods
	fmt.Println("Creating pod...")
	result, err := podsClient.Create(context.Background(), pods, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	if err != nil {
		log.Println("Error occcured while creating pod", err.Error())
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

}
func ListPods() {
	// List pods
	fmt.Printf("Listing pods in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}
}
func UpdatePods() {
	// Update services
	fmt.Println("Updating pods...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of service before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := podsClient.Get(context.Background(), "demo-pod", metav1.GetOptions{})
		if getErr != nil {
			//panic(fmt.Errorf("Failed to get latest version of service: %v", getErr))
			log.Println(fmt.Errorf("Failed to get latest version of pod: %v", getErr))
		}

		result.Spec.Containers[0].Image = "nginx:1.17" // reduce replica count
		_, updateErr := podsClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated pod...")
}
func DeletePods() {
	// Delete Pod
	fmt.Println("Deleting pod...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := podsClient.Delete(context.Background(), "demo-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted Pod.")

}
