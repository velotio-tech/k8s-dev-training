package serviceclient

import (
	"context"
	"fmt"

	util "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(svcName, namespace string) error {
	client := util.GetInClusterKubeConfigClient()
	svcClient := client.CoreV1().Services(namespace)
	svc := getServiceSpecs(svcName, namespace)

	_, err := svcClient.Create(context.Background(), svc, metav1.CreateOptions{})
	return err
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
	client := util.GetInClusterKubeConfigClient()
	svcClient := client.CoreV1().Services(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	return svcClient.Delete(context.Background(), svcName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

func UpdateService(svcName, namespace string) error {
	svcclient := util.GetInClusterKubeConfigClient().CoreV1().Services(namespace)
	svc, err := svcclient.Get(context.Background(), "my-svc", metav1.GetOptions{})
	if err != nil {
		return err
	}
	fmt.Println("svc generate name before update - ", svc.GetGenerateName())
	svc.SetGenerateName("svc-generate-name")
	_, err = svcclient.Update(context.Background(), svc, metav1.UpdateOptions{})
	fmt.Println("svc generate name after update - ", svc.GetGenerateName())
	return err
}

func GetAllServices(namespace string) (*corev1.ServiceList, error) {
	svcclient := util.GetInClusterKubeConfigClient().CoreV1().Services(namespace)
	return svcclient.List(context.Background(), metav1.ListOptions{})
}
