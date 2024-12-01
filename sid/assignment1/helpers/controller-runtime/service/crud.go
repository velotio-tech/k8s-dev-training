package service

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(client crtclient.Client, name string, labels map[string]string, nodePort, targetPort, port int) error {
	return client.Create(context.TODO(), &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceType("NodePort"),
			Ports: []corev1.ServicePort{
				{
					NodePort:   int32(nodePort),
					TargetPort: intstr.IntOrString{IntVal: int32(targetPort)},
					Port:       int32(port),
				},
			},
		},
	})
}

func Get(client crtclient.Client) error {
	serviceList := &corev1.ServiceList{}

	err := client.List(context.TODO(), serviceList)
	if err != nil {
		return err
	}

	for _, service := range serviceList.Items {
		fmt.Println(service.Name)
	}

	return nil
}

func Update(client crtclient.Client, name string, labels map[string]string) error {
	service := &corev1.Service{}

	err := client.Get(context.TODO(), crtclient.ObjectKey{Name: name}, service)
	if err != nil {
		return err
	}

	service.Labels = labels

	return client.Update(context.TODO(), service)
}

func Delete(client crtclient.Client, name string) error {
	return client.Delete(
		context.TODO(),
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: name}},
		&crtclient.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))},
	)
}
