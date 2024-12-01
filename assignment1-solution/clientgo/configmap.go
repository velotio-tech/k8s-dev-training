package clientgo

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateConfigmap(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("*** CREATE CONFIGMAP ***")
	new_cmap := &apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "configmap-demo",
		},
		Data: map[string]string{
			"ass-1": "configmap-lbl-demo",
		},
	}
	_, err := clientset.CoreV1().ConfigMaps("default").Create(ctx, new_cmap, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("Failed to create a configmap, error:", err)
		return
	}
	fmt.Println("Configmap created successfully!")
	fmt.Println("--------------------------------")
}

func ListConfigmaps(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("*** GET CONFIGMAP ***")
	cmap, err := clientset.CoreV1().ConfigMaps("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("Failed to get configmap, error:", err)
		return
	}
	fmt.Println("Total configmaps in the default namespace: ", len(cmap.Items))

	for _, item := range cmap.Items {
		fmt.Println(item.Name)
	}
	fmt.Println("--------------------------------")
}

func UpdateConfigmap(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("*** UPDATE CONFIGMAP ***")
	modified_cmap, err := clientset.CoreV1().ConfigMaps("default").Get(ctx, "configmap-demo", metav1.GetOptions{})
	if err != nil {
		fmt.Println("Failed to Get configmap, error:", err)
		return
	}

	if modified_cmap.Labels == nil {
		modified_cmap.Labels = make(map[string]string)
	}

	modified_cmap.Labels["app"] = "configmap-clientgo-lbl"
	_, err = clientset.CoreV1().ConfigMaps("default").Update(ctx, modified_cmap, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Failed to update the configmap, error: ", err)
		return
	}
	fmt.Println("Configmap updated successfully!")
	fmt.Println("--------------------------------")
}

func DeleteConfigmap(clientset *kubernetes.Clientset, ctx context.Context) {
	fmt.Println("*** DELETE CONFIGMAP ***")
	err := clientset.CoreV1().ConfigMaps("default").Delete(ctx, "configmap-demo", metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("Failed to delete configmaps, error:", err)
		return
	}
	fmt.Println("Configmap deleted successfully!")
	fmt.Println("--------------------------------")
}
