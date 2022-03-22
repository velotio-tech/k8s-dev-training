package controller

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateRtcPods(rtc client.Client) error {

	fmt.Println("Creating pod...")
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: "rtc-pod"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "rtc-pod" + "-container", Image: "nginx"}}},
	}
	err := rtc.Create(context.Background(), pod)
	if err != nil {
		return err
	}
	fmt.Println("Name : ", pod.Name, "Container Name : ", pod.Spec.Containers[0].Name, "Image : ", pod.Spec.Containers[0].Image)
	return nil
}

func ListRtcPods(rtc client.Client) error {
	pod := &corev1.PodList{}
	err := rtc.List(context.Background(), pod)
	if err != nil {
		return err
	} else {
		for _, p := range pod.Items {
			fmt.Println("Name : ", p.Name, "Labels : ", p.Labels)
		}
	}
	return nil
}

func UpdateRtcPods(rtc client.Client) error {

	fmt.Println("Updating pods...")
	pod := &corev1.Pod{}
	err := rtc.Get(context.TODO(), client.ObjectKey{
		Name: "rtc-pod",
	}, pod)
	if err != nil {
		return err
	}
	pod.Spec.Containers[0].Image = "nginx:1.17"
	err = rtc.Update(context.TODO(), pod)
	if err != nil {
		return err
	}
	fmt.Println(pod.Name, pod.Spec.Containers[0].Image)
	return nil
}

func DeleteRtcPods(rtc client.Client) error {

	fmt.Println("Deleting pod...")
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-pod",
		},
	}
	err := rtc.Delete(context.Background(), pod)
	if err != nil {
		return err
	}
	return nil
}
