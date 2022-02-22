package c_runtime

import (
	"assignment1/config"
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateService(name, namespace, svcType string, port int32) error {
	cl := config.GetClient()

	service := v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
			Namespace:    namespace,
			Labels:       map[string]string{"app": "deployment"},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{{
				Name:     "basic",
				Protocol: "TCP",
				Port:     port,
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: port,
				},
			}},
			Selector: map[string]string{"app": "deployment"},
			Type:     v1.ServiceType(svcType),
		},
		Status: v1.ServiceStatus{},
	}

	createOptions := client.CreateOptions{}

	err := cl.Create(context.Background(), &service, &createOptions)
	return err
}

func DeleteService(name, namespace string) error {

	cl := config.GetClient()

	deleteOptions := client.DeleteOptions{}
	service := v1.Service{}

	objectKey := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	cl.Get(context.Background(), objectKey, &service)

	err := cl.Delete(context.Background(), &service, &deleteOptions)

	return err
}

func ReadService(name, namespace string) error {

	var showServices bool
	if name == "" {
		showServices = true
	}

	cl := config.GetClient()

	if showServices {

		listOptions := client.ListOptions{}
		services := v1.ServiceList{}

		err := cl.List(context.Background(), &services, &listOptions)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(services, "", "\t")
		if err != nil {
			return err
		}

		log.Println(string(b))

	} else {

		getOptions := client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}
		service := v1.Service{}

		err := cl.Get(context.Background(), getOptions, &service)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(service, "", "\t")
		if err != nil {
			return err
		}

		log.Println(string(b))
	}

	return nil
}

func UpdateService(name, namespace string, port int32) error {

	cl := config.GetClient()

	getOptions := client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}
	service := v1.Service{}

	err := cl.Get(context.Background(), getOptions, &service)
	if err != nil {
		return err
	}

	if port != 0 {
		service.Spec.Ports[0].TargetPort.IntVal = port
	}

	updateOptions := client.UpdateOptions{}

	err = cl.Update(context.Background(), &service, &updateOptions)

	return err
}
