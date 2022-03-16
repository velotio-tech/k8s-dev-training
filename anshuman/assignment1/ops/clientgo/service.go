package clientgo

import (
	"assignment1/config"
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(name, namespace, svcType string, port int32) error {
	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

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

	createOptions := metav1.CreateOptions{}

	_, err := apiObj.Services(namespace).Create(context.Background(), &service, createOptions)
	return err
}

func DeleteService(name, namespace string) error {

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	deleteOptions := metav1.DeleteOptions{}

	err := apiObj.Services(namespace).Delete(context.Background(), name, deleteOptions)

	return err
}

func ReadService(name, namespace string) error {

	var showServices bool
	if name == "" {
		showServices = true
	}

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	if showServices {

		listOptions := metav1.ListOptions{}

		services, err := apiObj.Services(namespace).List(context.Background(), listOptions)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(services, "", "\t")
		if err != nil {
			return err
		}

		log.Println(string(b))

	} else {

		getOptions := metav1.GetOptions{}

		service, err := apiObj.Services(namespace).Get(context.Background(), name, getOptions)
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

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	getOptions := metav1.GetOptions{}

	service, err := apiObj.Services(namespace).Get(context.Background(), name, getOptions)
	if err != nil {
		return err
	}

	if port != 0 {
		service.Spec.Ports[0].TargetPort.IntVal = port
	}

	updateOptions := metav1.UpdateOptions{}

	_, err = apiObj.Services(namespace).Update(context.Background(), service, updateOptions)

	return err
}
