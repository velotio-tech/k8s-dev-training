/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package deployments

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// updateCmd represents the update deployment command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a deployment",
	Long:  "kube-client update deployment",
	Run: func(cmd *cobra.Command, args []string) {
		namespace := "default"
		deployment, err := GetDeployment()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if deployment.ObjectMeta.Labels == nil {
			deployment.ObjectMeta.Labels = make(map[string]string)
		}
		deployment.ObjectMeta.Labels["type"] = "frontend"

		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.Update(context.Background(), deployment)
		} else {
			_, err = cmdv1.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
		}
		if err != nil {
			fmt.Println("Failed to update deployment, error: ", err)
			return
		}

		fmt.Println("Deployment update successful")
	},
}

func init() {
	deploymentsCmd.AddCommand(updateCmd)
}
