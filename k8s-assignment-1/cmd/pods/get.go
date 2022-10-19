/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package pods

import (
	"context"
	"fmt"
	"os"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getCmd represents the get pod command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "kube-client pods get",
	Long:  "Get a pod",
	Run: func(cmd *cobra.Command, args []string) {
		pod, err := GetPod()
		if err != nil {
			fmt.Printf("failed to get pod in namespace default, error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Pod name: ", pod.Name)
	},
}

func init() {
	podsCmd.AddCommand(getCmd)
}

func GetPod() (*corev1.Pod, error) {
	var err error
	pod := &corev1.Pod{}
	namespace := "default"
	if cmdv1.UseCtrlRuntime {
		err = cmdv1.CtrlClient.Get(context.Background(), client.ObjectKey{
			Namespace: namespace,
			Name:      "my-pod",
		}, pod)
	} else {
		pod, err = cmdv1.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), "my-pod", v1.GetOptions{})
	}

	if err != nil {
		return nil, err
	}
	return pod, nil
}
