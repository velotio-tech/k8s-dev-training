package controllerruntime

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateConfigmap(clientset client.Client, ctx context.Context) {
	fmt.Println("*** CREATE CONFIGMAP ***")
	new_cmap := apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "configmap-demo",
		},
		Data: map[string]string{"ass-1": "controller-lbl-ass1"},
	}
	err := clientset.Create(ctx, &new_cmap)
	if err != nil {
		fmt.Println("Failed to create a configmap, error:", err)
		return
	}
	fmt.Println("Configmap created successfully!")
	fmt.Println("--------------------------------")
}

func ListConfigmaps(clientset client.Client, ctx context.Context) {
	fmt.Println("*** GET CONFIGMAP ***")
	cmap := &apiv1.ConfigMapList{}
	err := clientset.List(ctx, cmap)
	if err != nil {
		fmt.Println("Error during configmap list operation, error: ", err)
		return
	}
	fmt.Println("Total configmap: ", len(cmap.Items))
	for _, item := range cmap.Items {
		fmt.Println(item.Name)
	}
	fmt.Println("--------------------------------")
}

func UpdateConfigmap(clientset client.Client, ctx context.Context) {
	fmt.Println("*** UPDATE CONFIGMAP ***")
	cmap1 := &apiv1.ConfigMap{}
	err := clientset.Get(ctx, client.ObjectKey{Name: "demo-cm"}, cmap1)
	if err != nil {
		fmt.Println("Error during Get configmap operation, error:", err)
		return
	}

	if cmap1.Labels == nil {
		cmap1.Labels = make(map[string]string)
	}

	cmap1.Labels["app1"] = "controller-lbl-demo"
	err = clientset.Update(ctx, cmap1)
	if err != nil {
		fmt.Println("Failed to update configmap, error:", err)
		return
	}
	fmt.Println("Configmap updated successfully!")
	fmt.Println("--------------------------------")
}

func DeleteConfigmap(clientset client.Client, ctx context.Context) {
	fmt.Println("*** DELETE CONFIGMAP ***")
	del_cmap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "configmap-demo",
		},
	}
	err := clientset.Delete(ctx, del_cmap)
	if err != nil {
		fmt.Println("Delete cmap failed with an error : ", err)
	}
	fmt.Println("Configmap deleted successfully!")
	fmt.Println("--------------------------------")
}
