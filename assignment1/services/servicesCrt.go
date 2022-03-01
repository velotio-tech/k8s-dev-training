package services

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var err error
var rtc client.Client

func SetRtcClient(rtClient client.Client) {
	rtc = rtClient
}

func ListRtcServices() {

	fmt.Println("listing Services")
	services := &apiv1.ServiceList{}
	err = rtc.List(context.Background(), services)
	if err != nil {
		fmt.Println("error fetching pods list")
	} else {
		for _, each := range services.Items {
			fmt.Println("name: ", each.Name, " Labels: ", each.Labels)
		}
	}
}

func CreateRtcServices() {

	err = rtc.Create(context.Background(), service)
	if err != nil {
		fmt.Println(err)
	}
}

func UpdateRtcService() {

	svc := &apiv1.Service{}
	err = rtc.Get(context.Background(), client.ObjectKey{
		Name: "rtc-service",
	}, svc)
	if err != nil {
		fmt.Println(err)
	}
	svc.Spec.Ports[0].Protocol = "UDP"
	err = rtc.Update(context.Background(), svc)
}

func DeleteRtcService() {

	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-service",
		},
	}
	err = rtc.Delete(context.Background(), svc)
	if err != nil {
		fmt.Println(err)
	}
}
