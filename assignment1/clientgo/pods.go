package clientgo

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
)

var pods = &apiv1.Pod{
	TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
	ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
	Spec:       apiv1.PodSpec{Containers: []apiv1.Container{{Name: "test-nginx", Image: "nginx"}}},
}

func CreatePods(podsClient corev1.PodInterface) error {

	fmt.Println("Creating pod...")
	result, err := podsClient.Create(context.Background(), pods, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created pod %s.\n", result.GetObjectMeta().GetName())
	return nil
}

func ListAllPods(podsClient corev1.PodInterface) error {

	fmt.Printf("Listing pods in namespace %s:\n", apiv1.NamespaceDefault)
	list, err := podsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.Name)
	}
	return nil
}

func DeletePods(podsClient corev1.PodInterface) error {

	fmt.Println("Deleting pod...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := podsClient.Delete(context.Background(), "test-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	fmt.Println("Pod Deleted.")
	return nil
}

func UpdatePods(podsClient corev1.PodInterface) error {

	fmt.Println("Updating pods...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := podsClient.Get(context.Background(), "test-pod", metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}

		result.Spec.Containers[0].Image = "nginx:1.17"
		_, updateErr := podsClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return (fmt.Errorf("update Pod failed: %v", retryErr))
	}
	fmt.Println("Updated pod...")
	return nil
}
