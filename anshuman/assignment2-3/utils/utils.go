package utils

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodSpec(name, image string) v1.Pod {
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{{
				Name:  name,
				Image: image,
			}},
		},
	}

	return pod
}
