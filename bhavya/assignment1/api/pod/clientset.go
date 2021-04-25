package pod

import (
	"context"
	"fmt"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

func Get(client corev1.PodInterface, name string, options metav1.GetOptions) error {
	pod, err := client.Get(context.TODO(), name, options)
	if err != nil {
		return err
	}
	fmt.Println(pod.Name)
	return nil
}

func Create(client corev1.PodInterface, name, image string) error {
	pod := &v1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       v1.PodSpec{Containers: []v1.Container{{Name: name + "-container", Image: image}}},
	}
	_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func Update(client corev1.PodInterface, options metav1.UpdateOptions, name, image string) error {
	pod, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	pod.Spec.Containers[0].Image = image
	_, err = client.Update(context.TODO(), pod, options)
	if err != nil {
		return err
	}
	return nil
}

func Delete(client corev1.PodInterface, name string, options metav1.DeleteOptions) error {
	return client.Delete(context.TODO(), name, options)
}

func Operations() error {
	ns := "default"
	client := common.PodClient(ns)
	err := Create(client, "testpod", "nginx")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testpod", metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = Update(client, metav1.UpdateOptions{}, "testpod", "nginx:1.17" )
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testpod", metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = Delete(client, "testpod", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testpod", metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}