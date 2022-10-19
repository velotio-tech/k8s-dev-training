/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package pods

import (
	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
)

// podsCmd represents the pods command
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Perform operations on pods",
	Long:  "Perform operations on pods",
}

func init() {
	cmdv1.RootCmd.AddCommand(podsCmd)
}
