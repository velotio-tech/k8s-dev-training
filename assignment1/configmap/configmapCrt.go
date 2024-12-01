package configmap

import (
	"context"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var err error
var rtc client.Client

func SetRtClient (rtClient client.Client) {
	rtc = rtClient
}

func ListRTCConfigMaps(){
	fmt.Println("listing pods")
	configMaps := &apiv1.ConfigMapList{}
	err = rtc.List(context.Background(), configMaps)
	if err != nil {
		fmt.Println("error fetching pod list")
	} else {
		for _, each := range configMaps.Items {
			fmt.Println("name: ", each.Name, " Lables: ", each.Labels)
		}
	}
}

func CreateRTCConfigMap() {
	configMapData := make(map[string]string, 0)
	uiProperties := `allow.textmode=true`
	configMapData["ui.properties"] = uiProperties
	cm := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default-configmap",
		},
		Data: configMapData,
	}
	err = rtc.Create(context.Background(),cm)
	if err != nil {
		fmt.Println(err)
	}
}

func UpdateRTCConfigMap() {
	cm := &apiv1.ConfigMap{}
	configMapData := make(map[string]string, 0)
	uiProperties := `allow.textmode=false`
	configMapData["ui.properties"] = uiProperties

	err = rtc.Get(context.Background(),client.ObjectKey{
		Name: "default-configmap",
	}, cm)
	if err != nil {
		fmt.Println( err)
	}
	cm.Data = configMapData
	err = rtc.Update(context.Background(),cm)
}

func DeleteRTCConfigMap() {
	cm := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default-configmap",
		},
	}
	err = rtc.Delete(context.Background(),cm)
	if err != nil {
		fmt.Println(err)
	}
}

