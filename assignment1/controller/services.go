package controller

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var err error

func ListRtcServices(rtc client.Client) error {

	fmt.Println("listing Services")
	services := &apiv1.ServiceList{}
	err = rtc.List(context.Background(), services)
	if err != nil {
		return fmt.Errorf("error fetching pods list")
	} else {
		for _, each := range services.Items {
			fmt.Println("name: ", each.Name, " Labels: ", each.Labels)
		}
	}
	return nil
}

func CreateRtcServices(rtc client.Client) error {

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-service",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "demo",
			},
			Ports: []apiv1.ServicePort{
				{
					Name:     "access-port",
					Protocol: "TCP",
					Port:     8010,
				},
			},
		},
	}
	err = rtc.Create(context.Background(), service)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRtcService(rtc client.Client) error {

	svc := &apiv1.Service{}
	err = rtc.Get(context.Background(), client.ObjectKey{
		Name: "rtc-service",
	}, svc)
	if err != nil {
		return err
	}
	svc.Spec.Ports[0].Protocol = "UDP"
	err = rtc.Update(context.Background(), svc)
	return err
}

func DeleteRtcService(rtc client.Client) error {

	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rtc-service",
		},
	}
	err = rtc.Delete(context.Background(), svc)
	if err != nil {
		return err
	}
	return nil
}
