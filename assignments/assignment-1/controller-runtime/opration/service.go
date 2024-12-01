package opration

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createService() {

	createservice := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				labelkey: labelvalue,
			},
			Ports: []corev1.ServicePort{
				{
					Port: port,
				},
			},
		},
	}

	err := clientset.Create(context.Background(), createservice)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("service created")
}

func readService() {

	services := &v1.ServiceList{}
	err := clientset.List(context.Background(), services)
	if err != nil {
		log.Println(err.Error())
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	format := "%v\t%v\n"

	fmt.Fprintf(writer, format, "NAME")
	for _, service := range services.Items {

		fmt.Fprintf(writer, format, service.Name)
	}

	writer.Flush()
}

func updateService() {
	createservice := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				labelkey: labelvalue,
				"app":    "fronend",
			},
			Ports: []corev1.ServicePort{
				{
					Port: port,
				},
			},
		},
	}

	err := clientset.Update(context.Background(), createservice)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("service updated")
}

func deleteService() {

	err := clientset.Delete(context.Background(), &metav1.PartialObjectMetadata{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("service deleted")
}
