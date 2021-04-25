package configmap

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
	cm := &corev1.ConfigMap{}

	err := clt.Get(context.Background(),client.ObjectKey{
		Name: name,
	},cm)
	if err != nil {
		return err
	}
	fmt.Println(cm)
	return nil
}

func RtcCreate(clt client.Client,name string,data map[string]string) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Data: data,
	}
	err := clt.Create(context.Background(),cm)
	if err != nil {
		return err
	}
	return nil
}

func RtcUpdate(clt client.Client,name string,data map[string]string) error {
	cm := &corev1.ConfigMap{}

	err := clt.Get(context.Background(),client.ObjectKey{
		Name: name,
	}, cm)
	if err != nil {
		return err
	}
	cm.Data = data
	err = clt.Update(context.Background(),cm)
	fmt.Println(cm.Data)
	return nil
}

func RtcDelete(clt client.Client,name string) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := clt.Delete(context.Background(),cm)
	if err != nil {
		return err
	}
	return nil
}

func RtcOperation() error {
	ns := "default"
	data := make(map[string]string)
	data["key"] = "value"
	client := common.RcontrollerClient(ns)
	err := RtcCreate(client, "testcm", data)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testcm")
	if err != nil {
		return err
	}
	data["key"] = "value1"
	err = RtcUpdate(client, "testcm", data)
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testcm")
	if err != nil {
		return err
	}
	err = RtcDelete(client, "testcm")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testcm")
	if err != nil {
		return err
	}
	return nil
}