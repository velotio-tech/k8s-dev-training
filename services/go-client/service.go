package goclient

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Service struct {
	client *GoClient
}

func (s *Service) Create(ctx context.Context) error {
	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "service-1",
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
	options := metav1.CreateOptions{}
	_, err := s.client.CoreV1().Services("namespace-1").Create(ctx, svc, options)
	return err
}

func (s *Service) List(ctx context.Context) error {
	options := metav1.ListOptions{}
	svcs, err := s.client.CoreV1().Services("namespace-1").List(ctx, options)
	if err != nil {
		return err
	}
	for _, svc := range svcs.Items {
		fmt.Println(svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, svc.Spec.Ports, time.Since(svc.CreationTimestamp.Time).String())
	}
	return nil
}

func (s *Service) Update(ctx context.Context) error {
	svc, err := s.client.CoreV1().Services("namespace-1").Get(ctx, "service-1", metav1.GetOptions{})
	if err != nil {
		return err
	}
	if svc.Spec.Selector == nil {
		svc.Spec.Selector = make(map[string]string)
	}
	svc.Spec.Selector["label1-key"] = "label1-value"
	options := metav1.UpdateOptions{}
	_, err = s.client.CoreV1().Services("namespace-1").Update(ctx, svc, options)
	return err
}

func (s *Service) Delete(ctx context.Context) error {
	options := metav1.NewDeleteOptions(0)
	err := s.client.CoreV1().Services("namespace-1").Delete(ctx, "service-1", *options)
	return err
}
