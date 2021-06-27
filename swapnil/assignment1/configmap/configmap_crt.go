package main

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helpers "github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers"
	crcli "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	cli := helpers.GetCRTClient(metav1.NamespaceDefault)

	// creating pod
	helpers.Prompt()
	fmt.Println("creating configmap")
	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-configmap",
		},
		Data: map[string]string{"NAME":"Swapnil", "MOBILE": "8888888888"},
	}

	err := cli.Create(context.TODO(), configMap)
	if err != nil {
		fmt.Println("Error creating the configmap", err)
	} else {
		// check if the pod is created
		configMap := &apiv1.ConfigMap{}
		err := cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-configmap"}, configMap)
		//err := cli.List(context.TODO(), podList)
		if err == nil {
			fmt.Println(configMap.Name, configMap.Data)
		} else {
			fmt.Println("Error", err)
		}

		fmt.Println("config map creted")
	}
	helpers.Prompt()
	fmt.Println("updating config map")
	configMap = &apiv1.ConfigMap{}
	err = cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-configmap"}, configMap)
	if err != nil {
		fmt.Println("Error fetching config map", err)
	}else{
		fmt.Println("updating the config map")
		configMap.Data["APPNAME"] = "Awesome"
		err := cli.Update(context.TODO(), configMap)
		if err != nil {
			fmt.Println("Error updating the config map", err)
		} else {
			// verify if the lables are added
			configMap = &apiv1.ConfigMap{}
			err := cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-configmap"}, configMap)
			//err := cli.List(context.TODO(), podList)
			if err == nil {
				fmt.Println(configMap.Name, configMap.Data)
			} else {
				fmt.Println("Error", err)
			}
		}
	}

	// listing config maps
	helpers.Prompt()
	configMapList := &apiv1.ConfigMapList{}
	err = cli.List(context.TODO(), configMapList)
	if err != nil {
		fmt.Println("error fetching pod list")
	} else {
		for _, each := range configMapList.Items {
			fmt.Println("name: ", each.Name, " Lables: ", each.Data)
		}
	}

	// deleting the pod
	helpers.Prompt()
	configMap = &apiv1.ConfigMap{}
	err = cli.Get(context.TODO(), crcli.ObjectKey{Name:"demo-configmap"}, configMap)
	if err != nil {
		fmt.Println("Error fetching config map", err)
	} else {
		cli.Delete(
			context.TODO(),
			configMap,
			&crcli.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))})

		fmt.Println("config mpa deleted")
	}

}

