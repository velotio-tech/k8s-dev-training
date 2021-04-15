package main
 import (
	// "bufio"
	 "fmt"
	"context"
	
	// appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	//"k8s.io/kubernetes/pkg/api/unversioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helpers "github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers"

 )

func main() {
	clientset := helpers.GetClientset()
	podsClient := clientset.CoreV1().Pods(apiv1.NamespaceDefault)
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-pod",
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name: "nginx",
					Image: "nginx",
				},
			},
		},
	}
	// creating pod
	fmt.Println("Enter to create pod")
	helpers.Prompt()
	fmt.Println("Creating pod...")
	result, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created pod %q.\n", result.GetObjectMeta().GetName())

	fmt.Println("Enter to update pod")
	helpers.Prompt()
	fmt.Println("updating pod")
	thePod, getErr := podsClient.Get(context.TODO(), "demo-pod", metav1.GetOptions{})
	if getErr != nil {
		panic("couldn't get pod details")
	}
	thePod.Spec.Containers[0].Image = "nginx:1.13"
	_, updateErr := podsClient.Update(context.TODO(), thePod, metav1.UpdateOptions{})
	if updateErr != nil {
		panic("couldn't update pod")
	}
	fmt.Println("updated pod...")

	fmt.Println("Enter to get pods")
	helpers.Prompt()
	fmt.Println("getting pods")
	pods, listPodsErr := podsClient.List(context.TODO(), metav1.ListOptions{})
	if listPodsErr != nil {
		panic("couldn't get list of pods")
	}
	for _, d := range pods.Items {
		fmt.Printf(" * pod name - %s \n", d.Name)
	}
	fmt.Printf("pods fetched\n")

	fmt.Println("Enter to delete pod")
	helpers.Prompt()
	fmt.Println("deleting pod")
	deletePolicy := metav1.DeletePropagationForeground
	if err := podsClient.Delete(context.TODO(), "demo-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("pod deleted\n")
}