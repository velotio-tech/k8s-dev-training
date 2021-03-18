package pod

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Get(client crtclient.Client) error {
	podList := &corev1.PodList{}
	err := client.List(context.TODO(), podList)

	if err != nil {
		fmt.Println("Failed to get podlist cuz", err)
	}

	for _, pod := range podList.Items {
		fmt.Println(pod.Name)
	}
	return nil
}

func Create(client crtclient.Client, name, image string) error {
	return client.Create(context.TODO(), &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: name + "-container", Image: image}}},
	})
}

func Update(client crtclient.Client, name string, labels map[string]string) error {
	pod := &corev1.Pod{}

	err := client.Get(context.TODO(), crtclient.ObjectKey{Name: name}, pod)
	if err != nil {
		return err
	}

	pod.Labels = labels

	return client.Update(context.TODO(), pod)
}

func Delete(client crtclient.Client, name string) error {
	return client.Delete(
		context.TODO(),
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: name}},
		&crtclient.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))},
	)
}
