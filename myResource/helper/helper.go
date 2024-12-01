package helper

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilv2 "myResource/api/v2"
	"strings"
	"time"
)

func TimeLeft(executionTime string) (time.Duration, error) {
	now := time.Now().UTC()
	layout := "2006-01-02T15:04:05Z"
	scheduledTime, err := time.Parse(layout, executionTime)
	if err != nil {
		return time.Duration(0), err
	}
	return scheduledTime.Sub(now), nil
}

func SetPod(cr *utilv2.MyResource) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": cr.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: strings.Split(cr.Spec.Command, " "),
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}
