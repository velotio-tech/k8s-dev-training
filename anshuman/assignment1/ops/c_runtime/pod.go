package c_runtime

import (
	"assignment1/config"
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreatePod(name, namespace, image string) error {
	podDef := v1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
			Namespace:    namespace,
		},
		Spec: v1.PodSpec{
			Volumes:        []v1.Volume{},
			InitContainers: []v1.Container{},
			Containers: []v1.Container{{
				Name:  name,
				Image: image,
			}},
		},
	}

	createOptions := client.CreateOptions{
		Raw: &metav1.CreateOptions{},
	}

	cl := config.GetClient()

	err := cl.Create(context.Background(), &podDef, &createOptions)

	return err
}

func DeletePod(name, namespace string) error {

	deleteOptions := client.DeleteOptions{}

	cl := config.GetClient()

	pod := v1.Pod{}
	objectKey := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	cl.Get(context.Background(), objectKey, &pod)

	err := cl.Delete(context.Background(), &pod, &deleteOptions)

	return err
}

func ReadPod(name, namespace string) error {

	showPods := false
	if name == "" {
		showPods = true
	}

	cl := config.GetClient()

	if showPods {

		listOptions := client.ListOptions{
			Namespace: namespace,
		}

		pods := v1.PodList{}

		err := cl.List(context.Background(), &pods, &listOptions)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(pods, "", "\t")
		if err != nil {
			return err
		}
		log.Println(string(b))
	} else {

		pod := v1.Pod{}
		getOptions := client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}

		err := cl.Get(context.Background(), getOptions, &pod)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(pod, "", "\t")
		if err != nil {
			return err
		}
		log.Println(string(b))
	}
	return nil
}

func UpdatePod(name, namespace, image string) error {

	cl := config.GetClient()

	getOptions := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
	pod := v1.Pod{}
	updateOptions := client.UpdateOptions{}

	err := cl.Get(context.Background(), getOptions, &pod)
	if err != nil {
		return err
	}

	pod.Spec.Containers[0].Image = image

	err = cl.Update(context.Background(), &pod, &updateOptions)

	return err
}
