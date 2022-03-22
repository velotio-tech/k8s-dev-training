package controllers

import (
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func getDeploymentSpec(name, namespace, image, ownerName, ownerUID, ownerAPIVersion string, replicas *int32) appv1.Deployment {
	flag := true
	deploy := appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: ownerAPIVersion,
				Kind:       "Scaler",
				Name:       ownerName,
				UID:        types.UID(ownerUID),
				Controller: &flag,
			}},
			Namespace: namespace,
			Labels:    map[string]string{"app": "deployment"},
		},
		Spec: appv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "deployment"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:         name,
					GenerateName: name,
					Namespace:    name,
					Labels:       map[string]string{"app": "deployment"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  name,
						Image: image,
					}},
				},
			},
		},
	}
	return deploy
}
