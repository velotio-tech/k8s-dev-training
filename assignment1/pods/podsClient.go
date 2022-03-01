package pods

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

var pods = &apiv1.Pod{
	TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
	ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
	Spec:       apiv1.PodSpec{Containers: []apiv1.Container{{Name: "test-nginx", Image: "nginx"}}},
}

var podsClient corev1.PodInterface

func SetPodsClient(clientSet *kubernetes.Clientset) {
	podsClient = clientSet.CoreV1().Pods(apiv1.NamespaceDefault)
}

func CreatePods() {

	fmt.Println("Creating pod...")
	result, err := podsClient.Create(context.Background(), pods, metav1.CreateOptions{})
	if err != nil {
		log.Println("Error while creating pod", err.Error())
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

}

func ListAllPods() {

	fmt.Printf("Listing pods in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}
}

func DeletePods() {

	fmt.Println("Deleting pod...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := podsClient.Delete(context.Background(), "test-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Pod Deleted.")

}

func UpdatePods() {

	fmt.Println("Updating pods...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := podsClient.Get(context.Background(), "test-pod", metav1.GetOptions{})
		if getErr != nil {
			log.Println(fmt.Errorf("failed to get latest version of pod: %v", getErr))
		}

		result.Spec.Containers[0].Image = "nginx:1.17"
		_, updateErr := podsClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("update Pod failed: %v", retryErr))
	}
	fmt.Println("Updated pod...")
}
