/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package deployments

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteCmd represents the delete deployment command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a deployment",
	Long:  "kube-client delete deployment",
	Run: func(cmd *cobra.Command, args []string) {
		namespace := "default"
		var err error
		if cmdv1.UseCtrlRuntime {
			deployment, errGetDeployment := GetDeployment()
			if errGetDeployment != nil {
				fmt.Println("Error: ", errGetDeployment)
				return
			}
			err = cmdv1.CtrlClient.Delete(context.Background(), deployment)
		} else {
			err = cmdv1.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), "demo-deployment", v1.DeleteOptions{})
		}
		if err != nil {
			fmt.Println("Failed to delete deployment, error: ", err)
			return
		}

		fmt.Println("Deployment delete successful")
	},
}

func init() {
	deploymentsCmd.AddCommand(deleteCmd)
}
