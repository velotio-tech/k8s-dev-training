package podclient

import (
	"context"
	"fmt"

	constants "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/constants"
	util "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/utils"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespace string = core.NamespaceDefault
var podName string = "my-pod"

//	Public function to create pod.
func CreatePod(podName, ns string) error {
	client := util.GetInClusterKubeConfigClient()
	podClient := client.CoreV1().Pods(ns)
	pod := getPodSpecs(podName, ns)

	_, err := podClient.Create(context.Background(), pod, metav1.CreateOptions{})
	return err
}

//	Returns the specs of the pod to be created in cluster
func getPodSpecs(podName, ns string) *core.Pod {
	pod := &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "nginx-cont",
					Image:           constants.NginxImage,
					ImagePullPolicy: core.PullIfNotPresent,
				},
			},
		},
	}

	return pod
}

func DeletePod(podName string, ns string) error {
	client := util.GetInClusterKubeConfigClient()
	podClient := client.CoreV1().Pods(ns)
	deletePolicy := metav1.DeletePropagationForeground
	return podClient.Delete(context.TODO(), podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func UpdatePod(podName, ns string) error {
	// pod := &core.Pod{}
	apiClient := util.GetInClusterKubeConfigClient().CoreV1().Pods(core.NamespaceDefault)
	pod, err := apiClient.Get(context.Background(), "my-pod", metav1.GetOptions{})
	if err != nil {
		return err
	}

	fmt.Println("creation timestamp before update - ", pod.GetGenerateName())
	pod.SetGenerateName("generated-name")
	pod, err = apiClient.Update(context.Background(), pod, metav1.UpdateOptions{})
	fmt.Println("creation timestamp after update - ", pod.GetGenerateName())

	return err
}

func GetAllPods(namespace string) (*core.PodList, error) {
	podClient := util.GetInClusterKubeConfigClient().CoreV1().Pods(namespace)
	return podClient.List(context.Background(), metav1.ListOptions{})
}
