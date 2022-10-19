/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package namespace

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createCmd represents the create namespace command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create namespace",
	Long:  "kube-client create namespace",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		namespace := getNamespaceObj()

		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.Create(context.Background(), namespace)
		} else {
			_, err = cmdv1.ClientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
		}
		if err != nil {
			fmt.Println("Failed to create namespace. Error: ", err)
			return
		}

		fmt.Println("Namespace created successfully")
	},
}

func init() {
	namespaceCmd.AddCommand(createCmd)
}

func getNamespaceObj() *apiv1.Namespace {
	return &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-ns",
		},
	}
}
