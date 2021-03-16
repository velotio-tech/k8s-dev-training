package pod

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func GetPods(v1Client clientv1.CoreV1Interface, options metav1.ListOptions) error {
	pods, err := v1Client.Pods("default").List(context.TODO(), options)

	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	return nil
}

func DeletePods(v1Client clientv1.CoreV1Interface, listOptions metav1.ListOptions) error {
	return v1Client.Pods("default").DeleteCollection(
		context.TODO(), *metav1.NewDeleteOptions(int64(0)), listOptions)
}

func CreatePod(v1Client clientv1.CoreV1Interface, name, image string) error {
	containers := make([]corev1.Container, 1)
	containers[0] = corev1.Container{
		Name:  name + "-container",
		Image: image,
	}

	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{Containers: containers},
	}

	_, err := v1Client.Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})

	return err
}

func UpdatePods(v1Client clientv1.CoreV1Interface, options metav1.ListOptions, labels map[string]string) error {
	pods, err := v1Client.Pods("default").List(context.TODO(), options)

	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		pod.Labels = labels
		_, err := v1Client.Pods("default").Update(context.TODO(), &pod, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
