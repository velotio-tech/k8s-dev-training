package deployments

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var err error
var rtc client.Client

func SetRtClient (rtClient client.Client) {
	rtc = rtClient
}

func ListRTCDeployments(){
	fmt.Println("listing Deployments")
	deployments := &appsv1.DeploymentList{}
	err = rtc.List(context.Background(), deployments)
	if err != nil {
		fmt.Println("error fetching deployment list")
	} else {
		for _, each := range deployments.Items {
			fmt.Println("name: ", each.Name, " Lables: ", each.Labels)
		}
	}
}

func CreateRTCDeployment() {
	err = rtc.Create(context.Background(),deployment)
	if err != nil {
		fmt.Println(err)
	}
}

func UpdateRTCDeployment() {
	deploy := &appsv1.Deployment{}
	err = rtc.Get(context.Background(),client.ObjectKey{
		Name: "myService",
	}, deploy)
	if err != nil {
		fmt.Println( err)
	}
	deploy.Spec.Replicas = int32Ptr(1)
	err = rtc.Update(context.Background(), deploy)
}

func DeleteRTCDeployment() {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "myDeployment",
		},
	}
	err = rtc.Delete(context.Background(),deploy)
	if err != nil {
		fmt.Println(err)
	}
}

