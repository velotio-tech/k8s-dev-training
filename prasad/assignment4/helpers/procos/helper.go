package procos_helper

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "my.domain/ProcOS/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetProcOSDeploymentSpec(depName, namespace string, owner *batchv1.ProcOS) *appsv1.Deployment {
	var podTermnGracefulPeriod int64 = 0
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "procos-dep",
			Namespace:    namespace,
			Labels: map[string]string{
				"app": "procos-dep",
			},
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(owner, batchv1.GroupVersion.WithKind("ProcOS"))},
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "procos-dep",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "procos-dep",
					},
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: &podTermnGracefulPeriod,
					Containers: []corev1.Container{
						{
							Name:            "procos-job",
							Image:           "procos-job-image:latest",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"kubectl apply -f job.yaml",
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

func GetProcOSOwnedDeployments(ctx context.Context, cl client.Client, procos *batchv1.ProcOS) *appsv1.DeploymentList {
	depList := &appsv1.DeploymentList{}
	err := cl.List(ctx, depList, &client.ListOptions{Namespace: procos.Namespace})
	if err != nil {
		log.Println(err, ":: Unable to fetch deployments for procos owner...")
		return nil
	}

	return depList
}
