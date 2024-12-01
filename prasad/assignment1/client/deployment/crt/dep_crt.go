package dep_crt

import (
	"context"
	"fmt"

	util "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(depName, namespace string) error {
	crtClient, err := util.GetCRTClient()
	if err != nil {
		return err
	}

	return crtClient.Create(context.Background(), getDepSpecs(depName, namespace))
}

func getDepSpecs(depName, namespace string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      depName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "demo",
			},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "user-cont",
							Image:           "psdudeman39/user-ms:latest",
							ImagePullPolicy: corev1.PullNever,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8090,
								},
							},
						},
						{
							Name:            "order-cont",
							Image:           "psdudeman39/order-ms:latest",
							ImagePullPolicy: corev1.PullNever,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8091,
								},
							},
						},
						{
							Name:            "postgres-cont",
							Image:           "postgres:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Env: []corev1.EnvVar{
								{
									Name:  "POSTGRES_PASSWORD",
									Value: "password",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 5432,
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

func DeleteDeployment(depName, namespace string) error {
	client, err := util.GetCRTClient()
	if err != nil {
		return err
	}

	return client.Delete(context.Background(),
		&appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: depName, Namespace: namespace}},
		&crtclient.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))})
}

func UpdateDeployment(depName, namespace string) error {
	client, err := util.GetCRTClient()
	if err != nil {
		return err
	}
	dep := &appsv1.Deployment{}
	err = client.Get(context.TODO(), crtclient.ObjectKey{
		Name:      depName,
		Namespace: namespace,
	}, dep)
	if err != nil {
		return err
	}

	fmt.Println("Generate name before update - ", dep.GetGenerateName())
	dep.SetGenerateName("crt-generate-name")
	fmt.Println("Generate name after update - ", dep.GetGenerateName())
	return client.Update(context.Background(), dep)
}

func GetAllDeployments(namespace string) (*appsv1.DeploymentList, error) {
	client, err := util.GetCRTClient()
	if err != nil {
		return nil, err
	}
	depList := &appsv1.DeploymentList{}
	err = client.List(context.Background(), depList, &crtclient.ListOptions{})
	if err != nil {
		return nil, err
	}

	return depList, err
}

func int32Ptr(val int32) *int32 {
	return &val
}
