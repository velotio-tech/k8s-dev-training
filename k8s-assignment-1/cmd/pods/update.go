/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package pods

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// updateCmd represents the update pod command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a pod",
	Long:  "kube-client update pod",
	Run: func(cmd *cobra.Command, args []string) {
		namespace := "default"
		pod, err := GetPod()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if pod.ObjectMeta.Labels == nil {
			pod.ObjectMeta.Labels = make(map[string]string)
		}
		pod.ObjectMeta.Labels["type"] = "frontend"

		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.Update(context.Background(), pod)
		} else {
			_, err = cmdv1.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, v1.UpdateOptions{})
		}
		if err != nil {
			fmt.Println("Failed to update pod, error: ", err)
			return
		}

		fmt.Println("Pod update successful")
	},
}

func init() {
	podsCmd.AddCommand(updateCmd)
}
