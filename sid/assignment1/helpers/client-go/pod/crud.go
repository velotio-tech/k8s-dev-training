package pod

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func Get(client clientcorev1.PodInterface, options metav1.ListOptions) error {
	pods, err := client.List(context.TODO(), options)

	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	return nil
}

func Delete(client clientcorev1.PodInterface, listOptions metav1.ListOptions) error {
	return client.DeleteCollection(context.TODO(), *metav1.NewDeleteOptions(int64(0)), listOptions)
}

func Create(client clientcorev1.PodInterface, name, image string) error {
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: name + "-container", Image: image}}},
	}

	_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})

	return err
}

func Update(client clientcorev1.PodInterface, options metav1.ListOptions, labels map[string]string) error {
	pods, err := client.List(context.TODO(), options)

	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		pod.Labels = labels
		_, err := client.Update(context.TODO(), &pod, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
