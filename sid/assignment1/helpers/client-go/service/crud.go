package service

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func Get(client clientcorev1.ServiceInterface, options metav1.ListOptions) error {
	serviceList, err := client.List(context.TODO(), options)

	if err != nil {
		return err
	}

	for _, service := range serviceList.Items {
		fmt.Println(service.Name)
	}

	return nil
}

func Create(
	client clientcorev1.ServiceInterface,
	name string,
	labels map[string]string,
	nodePort, targetPort, port int,
) error {
	service := &corev1.Service{
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
	}

	_, err := client.Create(context.TODO(), service, metav1.CreateOptions{})

	return err
}

func Update(client clientcorev1.ServiceInterface, name string, labels map[string]string) error {
	service, err := client.Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		return err
	}

	service.Labels = labels
	_, err = client.Update(context.TODO(), service, metav1.UpdateOptions{})

	return err
}

func Delete(client clientcorev1.ServiceInterface, name string) error {
	return client.Delete(context.TODO(), name, *metav1.NewDeleteOptions(int64(0)))
}
