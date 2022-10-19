/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package namespace

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// updateCmd represents the update namespace command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a namespace",
	Long:  "kube-client update namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := GetNamespace()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if namespace.ObjectMeta.Labels == nil {
			namespace.ObjectMeta.Labels = make(map[string]string)
		}
		namespace.ObjectMeta.Labels["type"] = "frontend"

		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.Update(context.Background(), namespace)
		} else {
			_, err = cmdv1.ClientSet.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
		}
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Println("Namespace update successful")
	},
}

func init() {
	namespaceCmd.AddCommand(updateCmd)
}
