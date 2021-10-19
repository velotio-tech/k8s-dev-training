package service_crt

import (
	"context"
	"fmt"

	util "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateService(svcName, namespace string) error {
	crtClient, err := util.GetCRTClient()
	if err != nil {
		return err
	}
	return crtClient.Create(context.Background(), getServiceSpecs(svcName, namespace))
}

func getServiceSpecs(svcName, namespace string) *corev1.Service {
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       8091,
					TargetPort: intstr.FromInt(8091),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	return svc
}

func DeleteService(svcName, namespace string) error {
	client, err := util.GetCRTClient()
	if err != nil {
		return err
	}

	return client.Delete(context.Background(),
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: svcName, Namespace: namespace}},
		&crtclient.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))})
}

func GetAllServices(namespace string) (*corev1.ServiceList, error) {
	client, err := util.GetCRTClient()
	if err != nil {
		return nil, err
	}
	svcList := &corev1.ServiceList{}
	err = client.List(context.Background(), svcList, &crtclient.ListOptions{})
	if err != nil {
		return nil, err
	}

	return svcList, err
}

func UpdateService(svcName, namespace string) error {
	client, err := util.GetCRTClient()
	if err != nil {
		return err
	}
	svc := &corev1.Service{}
	err = client.Get(context.Background(), crtclient.ObjectKey{
		Name:      svcName,
		Namespace: namespace,
	}, svc)
	if err != nil {
		return err
	}

	fmt.Println("Generate name before update - ", svc.GetGenerateName())
	svc.SetGenerateName("crt-generate-name")
	fmt.Println("Generate name after update - ", svc.GetGenerateName())
	return client.Update(context.Background(), svc)
}
