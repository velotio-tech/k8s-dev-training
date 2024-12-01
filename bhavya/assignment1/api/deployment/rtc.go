package deployment

import (
	"context"
	"fmt"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/common"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func RtcGet(clt client.Client,name string) error {
	deploy := &appsv1.Deployment{}

	err := clt.Get(context.Background(),client.ObjectKey{
		Name: name,
	},deploy)
	if err != nil {
		return err
	}
	fmt.Println(deploy)
	return nil
}

func RtcCreate(clt client.Client,name,image string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []v1.ContainerPort{
								{
									Name:          "port1",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	err := clt.Create(context.Background(),deployment)
	if err != nil {
		return err
	}
	fmt.Println(deployment.Name,deployment.Spec.Template.Spec.Containers[0].Image)
	return nil
}

func RtcUpdate(clt client.Client,name,image string,) error {
	deployment := &appsv1.Deployment{}
	err := clt.Get(context.TODO(), client.ObjectKey{
		Name: name,
	},deployment)
	if err != nil {
		return err
	}
	deployment.Spec.Template.Spec.Containers[0].Image = image
	err = clt.Update(context.TODO(),deployment)
	if err != nil {
		return err
	}
	fmt.Println(deployment.Name,deployment.Spec.Template.Spec.Containers[0].Image)
	return nil
}

func RtcDelete(clt client.Client,name string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := clt.Delete(context.Background(),deployment)
	if err != nil {
		return err
	}
	return nil
}

func RtcOperation() error {
	ns := "default"
	client := common.RcontrollerClient(ns)
	err := RtcCreate(client, "testdeploy","nginx")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testdeploy")
	if err != nil {
		return err
	}
	err = RtcUpdate(client, "testdeoloy", "nginx:1.17")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testdeploy")
	if err != nil {
		return err
	}
	err = RtcDelete(client, "testdeploy")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testdeploy")
	if err != nil {
		return err
	}
	return nil
}