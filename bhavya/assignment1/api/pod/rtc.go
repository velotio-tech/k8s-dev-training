package pod

import (
	"context"
	"fmt"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func RtcGet(clt client.Client,name string) error {
	pod := &corev1.Pod{}

	err := clt.Get(context.Background(),client.ObjectKey{
		Name: name,
	},pod)
	if err != nil {
		return err
	}
	fmt.Println(pod)
	return nil
}

func RtcCreate(clt client.Client,name,image string) error {
	pod := &corev1.Pod{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Pod"},
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: name + "-container", Image: image}}},
	}
	err := clt.Create(context.Background(),pod)
	if err != nil {
		return err
	}
	return nil
}

func RtcUpdate(clt client.Client,name,image string,) error {
	pod := &corev1.Pod{}
	err := clt.Get(context.TODO(), client.ObjectKey{
		Name: name,
	},pod)
	if err != nil {
		return err
	}
	pod.Spec.Containers[0].Image = image
	err = clt.Update(context.TODO(),pod)
	if err != nil {
		return err
	}
	fmt.Println(pod.Name,pod.Spec.Containers[0].Image)
	return nil
}

func RtcDelete(clt client.Client,name string) error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := clt.Delete(context.Background(),pod)
	if err != nil {
		return err
	}
	return nil
}

func RtcOperation() error {
	ns := "default"
	client := common.RcontrollerClient(ns)
	err := RtcCreate(client, "testpod", "nginx")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testpod")
	if err != nil {
		return err
	}
	err = RtcUpdate(client, "testpod" , "nginx:1.17" )
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testpod")
	if err != nil {
		return err
	}
	err = RtcDelete(client, "testpod")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testpod")
	if err != nil {
		return err
	}
	return nil
}