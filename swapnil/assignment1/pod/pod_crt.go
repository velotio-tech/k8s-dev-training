package main

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helpers "github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers"
	crcli "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	cli := helpers.GetCRTClient(metav1.NamespaceDefault)

	// creating pod
	fmt.Println("creating pod")
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

	err := cli.Create(context.TODO(), pod)
	if err != nil {
		fmt.Println("Error creating the pod", err)
	} else {
		// check if the pod is created
		pod := &apiv1.Pod{}
		err := cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-pod"}, pod)
		//err := cli.List(context.TODO(), podList)
		if err == nil {
			fmt.Println(pod.Name, pod.Labels)
		} else {
			fmt.Println("Error", err)
		}

		fmt.Println("Pod creted")
	}
	fmt.Println("updating pod")
	pod_update := &apiv1.Pod{}
	err = cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-pod"}, pod_update)
	if err != nil {
		fmt.Println("Error fetching pod", err)
	}else{
		fmt.Println("updating the pod")
		pod_update.Labels = map[string]string{"APPNAME": "Awesome"}
		err := cli.Update(context.TODO(), pod_update)
		if err != nil {
			fmt.Println("Error updating the pod", err)
		} else {
			// verify if the lables are added
			pod = &apiv1.Pod{}
			err := cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-pod"}, pod)
			//err := cli.List(context.TODO(), podList)
			if err == nil {
				fmt.Println(pod.Name, pod.Labels)
			} else {
				fmt.Println("Error", err)
			}
		}
	}

	// listing pods
	podList := &apiv1.PodList{}
	err = cli.List(context.TODO(), podList)
	if err != nil {
		fmt.Println("error fetching pod list")
	} else {
		for _, each := range podList.Items {
			fmt.Println("name: ", each.Name, " Lables: ", each.Labels)
		}
	}

	// deleting the pod
	pod = &apiv1.Pod{}
	err = cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-pod"}, pod)
	if err != nil {
		fmt.Println("Error fetching pod", err)
	} else {
		cli.Delete(
			context.TODO(),
			pod,
			&crcli.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))})
	}

}

