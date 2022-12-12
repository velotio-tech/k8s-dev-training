package opration

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	corev1 "k8s.io/api/core/v1"
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

	service, err := clientset.CoreV1().Services("default").Create(context.Background(), createservice, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("service", service.Name, "created")
}

func readService() {
	services, err := clientset.CoreV1().Services("default").List(context.Background(), metav1.ListOptions{})
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

	service, err := clientset.CoreV1().Services("default").Update(context.Background(), createservice, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("service", service.Name, "updated")
}

func deleteService() {

	err := clientset.CoreV1().Services("default").Delete(context.Background(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Println("service", "name", "deleted")
}
