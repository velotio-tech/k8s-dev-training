package goclient

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	client *GoClient
}

func (namespace *Namespace) Create(ctx context.Context) error {
	ns := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-1",
		},
	}
	options := metav1.CreateOptions{}
	_, err := namespace.client.CoreV1().Namespaces().Create(ctx, ns, options)
	return err
}

func (namespace *Namespace) List(ctx context.Context) error {
	namespaces, err := namespace.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	for _, ns := range namespaces.Items {
		fmt.Println(ns.Name, ns.Status.Phase, time.Since(ns.CreationTimestamp.Time).String(), ns.Labels)
	}
	return err
}

func (namespace *Namespace) Update(ctx context.Context) error {
	ns := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-1",
			Labels: map[string]string{
				"label1_key": "label1_value",
			},
		},
	}
	options := metav1.UpdateOptions{}
	_, err := namespace.client.CoreV1().Namespaces().Update(ctx, ns, options)
	return err
}

func (namespace *Namespace) Delete(ctx context.Context) error {
	options := metav1.NewDeleteOptions(0)
	err := namespace.client.CoreV1().Namespaces().Delete(ctx, "namespace-1", *options)
	return err
}
