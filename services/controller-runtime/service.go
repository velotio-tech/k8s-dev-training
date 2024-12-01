package controllerruntime

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service struct {
	client *ControllerRuntime
}

func (s *Service) Create(ctx context.Context) error {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-1",
			Namespace: "namespace-1",
		},
		Spec: v1.ServiceSpec{
			Type: "ClusterIP",
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}
	options := client.CreateOptions{}
	return s.client.Create(ctx, svc, &options)
}

func (s *Service) List(ctx context.Context) error {
	sl := v1.ServiceList{}
	err := s.client.List(ctx, &sl)
	if err != nil {
		return err
	}
	for _, svc := range sl.Items {
		fmt.Println(svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, svc.Spec.Ports, time.Since(svc.CreationTimestamp.Time).String())
	}
	return nil
}

func (s *Service) Update(ctx context.Context) error {
	svc := v1.Service{}
	err := s.client.Get(ctx, types.NamespacedName{Namespace: "namespace-1", Name: "service-1"}, &svc)
	if err != nil {
		return err
	}

	if svc.Spec.Selector == nil {
		svc.Spec.Selector = make(map[string]string)
	}
	svc.Spec.Selector["label1-key"] = "label1-value"
	return s.client.Update(ctx, &svc)
}

func (s *Service) Delete(ctx context.Context) error {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-1",
			Namespace: "namespace-1",
		},
	}
	return s.client.Delete(ctx, svc)
}
