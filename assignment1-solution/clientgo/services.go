package clientgo

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateService(clientset *kubernetes.Clientset, ctx context.Context) {

	fmt.Println("***CREATE SERVICE***")
	new_svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "svc-demo",
		},
		Spec: v1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(ctx, new_svc, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Error occured during service creation, error: ", err)
		return
	}
	fmt.Println("Service created successfully!")
	fmt.Println("--------------------------------")
}

func ListServices(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("***LIST SERVICE***")

	svc, err := clientset.CoreV1().Services("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error occured during get services: ", err)
		return
	}
	fmt.Println("Total services: ", len(svc.Items))

	for _, svc := range svc.Items {
		fmt.Println(svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, svc.Spec.Ports)
	}
	fmt.Println("--------------------------------")
}

func UpdateService(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("***UPDATE SERVICE***")
	svc1, err := clientset.CoreV1().Services("default").Get(ctx, "svc-demo", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Service not found, error:", err)
		return
	}

	if svc1.Spec.Selector == nil {
		svc1.Spec.Selector = make(map[string]string)
	}

	svc1.Spec.Selector["app-1"] = "client-go-app-svc"
	_, err = clientset.CoreV1().Services("default").Update(ctx, svc1, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Error occured during update service, error: ", err)
		return
	}
	fmt.Println("Service Updated successfully!")
	fmt.Println("--------------------------------")
}

func DeleteService(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("***DELETE SERVICE***")
	err := clientset.CoreV1().Services("default").Delete(ctx, "svc-demo", metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Error occured during delete operation!", err)
		return
	}
	fmt.Println("Service deleted successfully!")
	fmt.Println("--------------------------------")
}
