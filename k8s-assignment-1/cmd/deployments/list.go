/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package deployments

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// listCmd represents the list deployments command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "kube-client deployments list",
	Long:  "List deployments",
	Run: func(cmd *cobra.Command, args []string) {
		namespace := "default"
		deployments := &appsv1.DeploymentList{}
		var err error
		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.List(context.Background(), deployments, client.InNamespace(namespace))
		} else {
			deployments, err = cmdv1.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
		}

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Printf("Namespace: %s\n", namespace)
		fmt.Printf("Number of deployments in the cluster: %d\n", len(deployments.Items))
		for i, deployment := range deployments.Items {
			fmt.Printf("%d. %s\n", i+1, deployment.Name)
		}
	},
}

func init() {
	deploymentsCmd.AddCommand(listCmd)
}
