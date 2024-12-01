package controllerruntime

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateService(clientset client.Client, ctx context.Context) {
	fmt.Println("*** CREATE SERVICE ***")
	new_svc := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc-controller-demo",
		},
		Spec: apiv1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	err := clientset.Create(ctx, new_svc)
	if err != nil {
		fmt.Println("Service is not created, error: ", err)
		return
	}
	fmt.Println("Service created successfully!")
	fmt.Println("--------------------------------")
}

func ListServices(clientset client.Client, ctx context.Context) {
	fmt.Println("*** GET SERVICE ***")
	svc := &apiv1.ServiceList{}
	err := clientset.List(ctx, svc)
	if err != nil {
		fmt.Println("Error in list service operation, error: ", err)
		return
	}
	fmt.Println("Total services: ", len(svc.Items))

	for _, svc := range svc.Items {
		fmt.Println(svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, svc.Spec.Ports)
	}
	fmt.Println("--------------------------------")
}

func UpdateService(clientset client.Client, ctx context.Context) {
	fmt.Println("*** UPDATE SERVICE ***")
	svc1 := &apiv1.Service{}
	err := clientset.Get(ctx, client.ObjectKey{Name: "nginx-service"}, svc1)
	if err != nil {
		fmt.Println("Get service operation failed, error: ", err)
		return
	}

	if svc1.ObjectMeta.Labels == nil {
		svc1.ObjectMeta.Labels = make(map[string]string)
	}

	svc1.ObjectMeta.Labels["app1"] = "svc-controller-lbl1"
	err = clientset.Update(ctx, svc1)
	if err != nil {
		fmt.Println("Update service failed, error: ", err)
		return
	}
	fmt.Println("Service updated successfully!")
	fmt.Println("--------------------------------")
}

func DeleteService(clientset client.Client, ctx context.Context) {
	fmt.Println("*** DELETE SERVICE ***")
	del_svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc-controller-demo",
		},
	}

	err := clientset.Delete(ctx, del_svc)
	if err != nil {
		fmt.Println("Service is not deleted, error: ", err)
		return
	}
	fmt.Println("Service deleted successfully!")
	fmt.Println("--------------------------------")
}
