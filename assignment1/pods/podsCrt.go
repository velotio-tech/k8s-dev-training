package pods

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var rtc client.Client

func SetRtcClient(rtClient client.Client) {
	rtc = rtClient
}

func CreateRtcPods() {

	fmt.Println("Creating pod...")
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "test-pod" + "-container", Image: "nginx"}}},
	}
	err := rtc.Create(context.Background(), pod)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Name : ", pod.Name, "Container Name : ", pod.Spec.Containers[0].Name, "Image : ", pod.Spec.Containers[0].Image)
}

func ListRtcPods() {
	pod := &corev1.PodList{}
	err := rtc.List(context.Background(), pod)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, p := range pod.Items {
			fmt.Println("Name : ", p.Name, "Labels : ", p.Labels)
		}
	}

}

func UpdateRtcPods() {

	fmt.Println("Updating pods...")
	pod := &corev1.Pod{}
	err := rtc.Get(context.TODO(), client.ObjectKey{
		Name: "test-pod",
	}, pod)
	if err != nil {
		fmt.Println(err)
	}
	pod.Spec.Containers[0].Image = "nginx:1.17"
	err = rtc.Update(context.TODO(), pod)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pod.Name, pod.Spec.Containers[0].Image)
}

func DeleteRtcPods() {

	fmt.Println("Deleting pod...")
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
		},
	}
	err := rtc.Delete(context.Background(), pod)
	if err != nil {
		fmt.Println(err)
	}

}
