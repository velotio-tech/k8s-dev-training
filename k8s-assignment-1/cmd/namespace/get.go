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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// getCmd represents the get namespace command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "kube-client namespace get",
	Long:  "Get a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := GetNamespace()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Namespace name: ", namespace.Name)
	},
}

func init() {
	namespaceCmd.AddCommand(getCmd)
}

func GetNamespace() (*corev1.Namespace, error) {
	var err error
	namespace := &corev1.Namespace{}
	if cmdv1.UseCtrlRuntime {
		err = cmdv1.CtrlClient.Get(context.Background(), client.ObjectKey{
			Name: "demo-ns",
		}, namespace)
	} else {
		namespace, err = cmdv1.ClientSet.CoreV1().Namespaces().Get(context.TODO(), "demo-ns", v1.GetOptions{})
	}
	if err != nil {
		return nil, err
	}
	return namespace, nil
}
