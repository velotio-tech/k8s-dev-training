package configmap

import (
	"context"
	"fmt"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"time"
)

func Get(client corev1.ConfigMapInterface, name string, options metav1.GetOptions) error {
	cm, err := client.Get(context.TODO(), name, options)
	if err != nil {
		return err
	}
	fmt.Println(cm.Name, cm.Data)
	return nil
}

func Create(client corev1.ConfigMapInterface, name string, data map[string]string) error {
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	_, err := client.Create(context.TODO(), cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func Update(client corev1.ConfigMapInterface, options metav1.UpdateOptions, name string, data map[string]string) error {
	cm, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	cm.Data = data
	_, err = client.Update(context.TODO(), cm, options)
	if err != nil {
		return err
	}
	cm, err = client.Get(context.TODO(), name, metav1.GetOptions{})
	return nil
}

func Delete(client corev1.ConfigMapInterface, name string, options metav1.DeleteOptions) error {
	return client.Delete(context.TODO(), name, options)
}

func Operations() error {
	ns := "default"
	data := make(map[string]string)
	data["key"] = "value"
	client := common.ConfigMapClient(ns)
	err := Create(client, "testcm", data)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testcm", metav1.GetOptions{})
	if err != nil {
		return err
	}
	data["key"] = "value1"
	err = Update(client, metav1.UpdateOptions{}, "testcm", data)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testcm", metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = Delete(client, "testcm", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = Get(client, "testcm", metav1.GetOptions{})
	if err != nil {
		return err
	}
	return nil
}
