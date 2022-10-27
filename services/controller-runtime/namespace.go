package controllerruntime

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	client *ControllerRuntime
}

func (n *Namespace) Create(ctx context.Context) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-1",
		},
	}
	return n.client.Create(ctx, ns)
}

func (n *Namespace) List(ctx context.Context) error {
	nsList := v1.NamespaceList{}
	err := n.client.List(ctx, &nsList)
	if err != nil {
		return err
	}
	for _, ns := range nsList.Items {
		fmt.Println(ns.Name, ns.Status.Phase, time.Since(ns.CreationTimestamp.Time).String(), ns.Labels)
	}
	return nil
}

func (n *Namespace) Update(ctx context.Context) error {

	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-1",
			Labels: map[string]string{
				"label1_key": "label1_value",
			},
		},
	}
	return n.client.Update(ctx, ns)
}

func (n *Namespace) Delete(ctx context.Context) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "namespace-1",
		},
	}
	return n.client.Delete(ctx, ns)
}
