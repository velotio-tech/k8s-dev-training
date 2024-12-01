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

func createConfigMap() {

	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: configmapName,
		},
		Data: configData,
	}

	err := clientset.Create(context.Background(), configmap)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("configmap created")
}

func readConfigMap() {
	configs := &corev1.ComponentStatusList{}
	err := clientset.List(context.Background(), configs)
	if err != nil {
		log.Println(err.Error())
		return
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	format := "%v\t%v\n"

	fmt.Fprintf(writer, format, "NAME")
	for _, service := range configs.Items {

		fmt.Fprintf(writer, format, service.Name)
	}

	writer.Flush()
}

func updateConfigMap() {

	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: configmapName,
		},
		Data: configDataUpdated,
	}

	err := clientset.Update(context.Background(), configmap)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("configmap updated")
}

func deleteConfigMap() {

	err := clientset.Delete(context.Background(), &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name: configmapName,
		},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("configmap", configmapName, "deleted")

}
