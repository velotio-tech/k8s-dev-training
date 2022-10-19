/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package pods

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createCmd represents the create pod command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create pod",
	Long:  "kube-client create pod",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		namespace := "default"
		pod := getPodObj(namespace)
		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.Create(context.Background(), pod)
		} else {
			_, err = cmdv1.ClientSet.CoreV1().Pods(namespace).Create(context.TODO(), pod, v1.CreateOptions{})
		}
		if err != nil {
			fmt.Println("Error creating pod, error: ", err)
			return
		}
		fmt.Println("Pod created successfully")
	},
}

func getPodObj(namespace string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "my-pod",
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx-ap",
					Image: "nginx",
				},
			},
		},
	}
}

func init() {
	podsCmd.AddCommand(createCmd)
}
