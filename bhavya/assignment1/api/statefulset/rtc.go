package statefulset

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
	sfs := &appsv1.StatefulSet{}

	err := clt.Get(context.Background(),client.ObjectKey{
		Name: name,
	},sfs)
	if err != nil {
		return err
	}
	fmt.Println(sfs.Name,sfs.Spec.Template.Spec.Containers[0].Image)
	return nil
}

func RtcCreate(clt client.Client,name,image string) error {
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.StatefulSetSpec{
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
	err := clt.Create(context.Background(),statefulset)
	if err != nil {
		return err
	}
	return nil
}

func RtcUpdate(clt client.Client,name,image string,) error {
	sfs := &appsv1.StatefulSet{}
	err := clt.Get(context.TODO(), client.ObjectKey{
		Name: name,
	},sfs)
	if err != nil {
		return err
	}
	sfs.Spec.Template.Spec.Containers[0].Image = image
	err = clt.Update(context.TODO(),sfs)
	if err != nil {
		return err
	}
	fmt.Println(sfs.Name,sfs.Spec.Template.Spec.Containers[0].Image)
	return nil
}

func RtcDelete(clt client.Client,name string) error {
	sfs := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	err := clt.Delete(context.Background(),sfs)
	if err != nil {
		return err
	}
	return nil
}


func RtcOperation() error {
	ns := "default"
	client := common.RcontrollerClient(ns)
	err := RtcCreate(client, "testsfs", "nginx")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testsfs")
	if err != nil {
		return err
	}
	err = RtcUpdate(client, "testsfs", "nginx:1.17" )
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testsfs")
	if err != nil {
		return err
	}
	err = RtcDelete(client, "testsfs")
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)
	err = RtcGet(client, "testsfs")
	if err != nil {
		return err
	}
	return nil
}