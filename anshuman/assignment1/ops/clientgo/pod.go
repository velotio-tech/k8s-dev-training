package clientgo

import (
	"assignment1/config"
	"context"
	"encoding/json"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	createOptions := metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	_, err := apiObj.Pods(namespace).Create(context.Background(), &podDef, createOptions)

	return err
}

func DeletePod(name, namespace string) error {

	deleteOptions := metav1.DeleteOptions{}

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	err := apiObj.Pods(namespace).Delete(context.Background(), name, deleteOptions)

	return err
}

func ReadPod(name, namespace string) error {

	showPods := false
	if name == "" {
		showPods = true
	}

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	if showPods {

		listOptions := metav1.ListOptions{}

		pods, err := apiObj.Pods(namespace).List(context.Background(), listOptions)
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(pods, "", "\t")
		if err != nil {
			return err
		}
		log.Println(string(b))
	} else {

		getOptions := metav1.GetOptions{}

		pod, err := apiObj.Pods(namespace).Get(context.Background(), name, getOptions)
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

	cs := config.GetClientSet()
	apiObj := cs.CoreV1()

	getOptions := metav1.GetOptions{}
	updateOptions := metav1.UpdateOptions{}

	pod, err := apiObj.Pods(namespace).Get(context.Background(), name, getOptions)
	if err != nil {
		return err
	}

	pod.Spec.Containers[0].Image = image

	_, err = apiObj.Pods(namespace).Update(context.Background(), pod, updateOptions)

	return err
}
