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

	config, err := clientset.CoreV1().ConfigMaps("default").Create(context.Background(), configmap, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("configmap", config.Name, "created")
}

func readConfigMap() {
	configs, err := clientset.CoreV1().ConfigMaps("default").List(context.Background(), metav1.ListOptions{})
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

	config, err := clientset.CoreV1().ConfigMaps("default").Update(context.Background(), configmap, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("configmap", config.Name, "updated")
}

func deleteConfigMap() {

	err := clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), configmapName, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("configmap", configmapName, "deleted")

}
