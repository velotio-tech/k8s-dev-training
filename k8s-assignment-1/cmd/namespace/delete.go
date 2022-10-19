/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package namespace

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteCmd represents the delete namespace command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a namespace",
	Long:  "kube-client delete namespace",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if cmdv1.UseCtrlRuntime {
			ns, errGetNs := GetNamespace()
			if errGetNs != nil {
				fmt.Println("Error: ", errGetNs)
				return
			}
			err = cmdv1.CtrlClient.Delete(context.Background(), ns)
		} else {
			err = cmdv1.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), "demo-ns", v1.DeleteOptions{})
		}
		if err != nil {
			fmt.Println("Failed to delete namespace, error: ", err)
			return
		}

		fmt.Println("Namespace delete successful")
	},
}

func init() {
	namespaceCmd.AddCommand(deleteCmd)
}
