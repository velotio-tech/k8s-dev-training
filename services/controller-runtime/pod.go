package controllerruntime

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Pod struct {
	client *ControllerRuntime
}

func (p *Pod) Create(ctx context.Context) error {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-1",
			Namespace: "namespace-1",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-1",
					Image: "nginx:latest",
				},
			},
		},
	}
	options := client.CreateOptions{}
	return p.client.Create(ctx, pod, &options)

}

func (p *Pod) List(ctx context.Context) error {
	podList := v1.PodList{TypeMeta: metav1.TypeMeta{}}
	err := p.client.List(ctx, &podList)
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		fmt.Println(pod.Name, pod.Status.Phase, time.Since(pod.CreationTimestamp.Time).String(), pod.Labels)
	}
	return nil
}

func (p *Pod) Update(ctx context.Context) error {
	pod := v1.Pod{}
	err := p.client.Get(ctx, types.NamespacedName{Namespace: "namespace-1", Name: "pod-1"}, &pod)
	if err != nil {
		return err
	}
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}
	pod.Labels["label1-key"] = "label1-value"
	return p.client.Update(ctx, &pod)
}

func (p *Pod) Delete(ctx context.Context) error {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-1",
			Namespace: "namespace-1",
		},
	}
	return p.client.Delete(ctx, pod)
}
