package podcrt

import (
	"context"

	util "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crtclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func CreatePod(podName, namespace string) error {
	crtClient, err := util.GetCRTClient()
	if err != nil {
		return err
	}

	return crtClient.Create(context.Background(), getPodSpecs(podName, namespace))
}

func DeletePod(podName, namespace string) error {
	crtClient, err := util.GetCRTClient()
	if err != nil {
		return err
	}

	return crtClient.Delete(context.Background(),
		&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: podName, Namespace: namespace}},
		&crtclient.DeleteOptions{Raw: metav1.NewDeleteOptions(int64(0))})
}

func GetAllPods(namespace string) (*corev1.PodList, error) {
	client, err := util.GetCRTClient()
	if err != nil {
		return nil, err
	}
	podList := &corev1.PodList{}
	err = client.List(context.Background(), podList, &crtclient.ListOptions{})
	if err != nil {
		return nil, err
	}

	return podList, err
}

func UpdatePod(podName, namespace string) error {
	client, err := util.GetCRTClient()
	if err != nil {
		return err
	}
	pod := &corev1.Pod{}
	err = client.Get(context.Background(), crtclient.ObjectKey{
		Name:      podName,
		Namespace: namespace,
	}, pod)
	if err != nil {
		return err
	}

	pod.SetGenerateName("generate-name-edit")
	return client.Update(context.Background(), pod)
}

func getPodSpecs(podName, ns string) *corev1.Pod {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: ns,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Containers: []corev1.Container{
				{
					Name:            "nginx-cont-crt",
					Image:           "nginx:latest",
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
		},
	}

	return pod
}
