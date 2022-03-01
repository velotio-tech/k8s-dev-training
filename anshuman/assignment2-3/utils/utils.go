package utils

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func returnBoolPointer(b bool) *bool {
	return &b
}

func GetPodSpec(name, image, namespace, parentName, parentUID, parentAPIVersion string) v1.Pod {
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:         name,
			GenerateName: name,
			Namespace:    namespace,
			Labels:       map[string]string{"owner": parentName},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         parentAPIVersion,
				Kind:               "customScaler",
				Name:               parentName,
				UID:                types.UID(parentUID),
				Controller:         returnBoolPointer(true),
				BlockOwnerDeletion: new(bool),
			}},
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
