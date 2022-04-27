package controllerruntimecrud

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListPods(controllerClient client.Client) {
	podList := &corev1.PodList{}
	err := controllerClient.List(context.Background(), podList, client.InNamespace("default"))
	if err != nil {
		log.Printf("failed to list pods in namespace default: %v\n", err)
	} else {
		for _, p := range podList.Items {
			fmt.Println("Name : ", p.Name, "Labels : ", p.Labels)
		}
	}
}

func CreatePod(controllerClient client.Client) {
	// the yaml file content will come here
	newPod := &corev1.Pod{
		// parav this will be the same as client-go
		TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: "busybox-pod2",
			Labels: map[string]string{
				"owner": "parav",
			}},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "busybox", Image: "busybox:latest", Command: []string{"sleep", "10000"}},
			}},
	}
	err := controllerClient.Create(context.Background(), newPod)
	if err != nil {
		log.Printf("cannot create new pod using controller runtime: %v", err)
	} else {
		fmt.Println("POD Successfully Created.")
	}
}

func EditPod(controllerClient client.Client) {
	pod := &corev1.Pod{}

	//get the latest version of the desired pod
	err := controllerClient.Get(context.TODO(), client.ObjectKey{
		Name: "busybox-pod2",
	}, pod)
	if err != nil {
		log.Printf("cannot get current version of desired pod: %v", err)
	}
	//let us update the label of the pod
	pod.ObjectMeta.Labels["owner"] = "kaushal"
	err = controllerClient.Update(context.TODO(), pod)
	if err != nil {
		log.Printf("cannot update desired pod: %v", err)
	} else {
		fmt.Println("POD Successfully updated.")
	}
}

func DeletePod(controllerClient client.Client) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox-pod2",
		},
	}
	err := controllerClient.Delete(context.Background(), pod)
	if err != nil {
		log.Printf("cannot delete desired pod: %v", err)
	} else {
		fmt.Println("POD Successfully deleted.")
	}
}
