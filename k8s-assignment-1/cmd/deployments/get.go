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

// getCmd represents the get deployment command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "kube-client deployment get",
	Long:  "Get a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		deployment, err := GetDeployment()
		if err != nil {
			fmt.Println("failed to get deployment in namespace default, error: ", err)
			return
		}
		fmt.Println("Deployment name: ", deployment.Name)
	},
}

func init() {
	deploymentsCmd.AddCommand(getCmd)
}

func GetDeployment() (*appsv1.Deployment, error) {
	var err error
	deployment := &appsv1.Deployment{}
	namespace := "default"
	if cmdv1.UseCtrlRuntime {
		err = cmdv1.CtrlClient.Get(context.Background(), client.ObjectKey{
			Namespace: namespace,
			Name:      "demo-deployment",
		}, deployment)
	} else {
		deployment, err = cmdv1.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), "demo-deployment", v1.GetOptions{})
	}
	if err != nil {
		return nil, err
	}
	return deployment, nil
}
