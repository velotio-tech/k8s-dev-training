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

// listCmd represents the view command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "kube-client pods list",
	Long:  "List pods",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		pods := &corev1.PodList{}
		namespace := "default"
		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.List(context.Background(), pods, client.InNamespace(namespace))
		} else {
			pods, err = cmdv1.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
		}
		if err != nil {
			fmt.Printf("failed to list pods in namespace default: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Namespace: %s\n", namespace)
		fmt.Printf("Number of pods in the cluster: %d\n", len(pods.Items))

		for i, pod := range pods.Items {
			fmt.Printf("%d. %s\n", i+1, pod.Name)
		}

	},
}

func init() {
	podsCmd.AddCommand(listCmd)
}
