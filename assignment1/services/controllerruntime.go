package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var rtc client.Client

func CreateRTServiceClient(rtClient client.Client) {
	rtc = rtClient
}

func CreateRTServices() {
	fmt.Println("TODO")

}

func ListRTServices() {
	fmt.Println("listing Services")
	services := &corev1.ServiceList{}
	err := rtc.List(context.Background(), services)
	if err != nil {
		fmt.Println("error fetching pod list")
	} else {
		for _, each := range services.Items {
			fmt.Println("name: ", each.Name, " Lables: ", each.Labels)
		}
	}
}

func UpdateRTServices() {

	svc := &corev1.Service{}
	err := rtc.Get(context.Background(), client.ObjectKey{
		Name: "rtc-service",
	}, svc)

	if err != nil {
		fmt.Println(err)
	}
	svc.Spec.Ports[0].Protocol = "UDP"
	err = rtc.Update(context.Background(), svc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(svc.Name, svc.Spec.Ports[0].Protocol)

}

func DeleteRTServices() {
	services := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-service",
		},
	}
	err := rtc.Delete(context.Background(), services)
	if err != nil {
		fmt.Println(err)
	}

}
