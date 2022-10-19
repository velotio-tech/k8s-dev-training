/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package namespace

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// listCmd represents the list namespace command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "kube-client namespace list",
	Long:  "List namespace",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		namespaces := &corev1.NamespaceList{}
		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.List(context.Background(), namespaces)
		} else {
			namespaces, err = cmdv1.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		}
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Printf("Number of namespaces in the cluster: %d\n", len(namespaces.Items))
		for i, namespace := range namespaces.Items {
			fmt.Printf("%d. %s\n", i+1, namespace.Name)
		}
	},
}

func init() {
	namespaceCmd.AddCommand(listCmd)
}
