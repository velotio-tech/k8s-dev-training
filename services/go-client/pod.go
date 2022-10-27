package goclient

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	client *GoClient
}

func (p *Pod) Create(ctx context.Context) error {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-1",
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
	options := metav1.CreateOptions{}
	_, err := p.client.CoreV1().Pods("namespace-1").Create(ctx, pod, options)
	return err
}

func (p *Pod) List(ctx context.Context) error {
	options := metav1.ListOptions{}
	pods, err := p.client.CoreV1().Pods("namespace-1").List(ctx, options)
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Status.Phase, time.Since(pod.CreationTimestamp.Time).String(), pod.Labels)
	}
	return nil
}

func (p *Pod) Update(ctx context.Context) error {
	pod, err := p.client.CoreV1().Pods("namespace-1").Get(ctx, "pod-1", metav1.GetOptions{})
	if err != nil {
		return err
	}
	if pod.Labels == nil {
		pod.Labels = make(map[string]string)
	}
	pod.Labels["label1-key"] = "label1-value"
	options := metav1.UpdateOptions{}
	_, err = p.client.CoreV1().Pods("namespace-1").Update(ctx, pod, options)
	return err
}

func (p *Pod) Delete(ctx context.Context) error {
	options := metav1.NewDeleteOptions(0)
	err := p.client.CoreV1().Pods("namespace-1").Delete(ctx, "pod-1", *options)
	return err
}
