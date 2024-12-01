package depclient

import (
	"context"

	util "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateDeployment(depName, namespace string) error {
	client := util.GetInClusterKubeConfigClient()
	deploymentSpecs := getDeployment(depName, namespace)
	depClient := client.AppsV1().Deployments(namespace)
	_, err := depClient.Create(context.Background(), deploymentSpecs, metav1.CreateOptions{})
	return err
}

func DeleteDeployment(depName, namespace string) error {
	deletePolicy := metav1.DeletePropagationForeground
	depClient := util.GetInClusterKubeConfigClient().AppsV1().Deployments(namespace)
	err := depClient.Delete(context.Background(), depName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	return err
}

func GetAllDeployments(namespace string) (*appsv1.DeploymentList, error) {
	depClient := util.GetInClusterKubeConfigClient().AppsV1().Deployments(namespace)
	return depClient.List(context.Background(), metav1.ListOptions{})
}

func UpdateDeployment(depName, namespace string) error {
	depClient := util.GetInClusterKubeConfigClient().AppsV1().Deployments(namespace)
	deployment, err := depClient.Get(context.Background(), "inc-deployment", metav1.GetOptions{})
	if err != nil {
		return err
	}
	deployment.SetGenerateName("dep-generate-name")
	_, err = depClient.Update(context.Background(), deployment, metav1.UpdateOptions{})
	return err
}

func getDeployment(depName, namespace string) *appsv1.Deployment {
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

func int32Ptr(val int32) *int32 {
	return &val
}
