package namespace

import (
	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
)

// namespaceCmd represents the namespace command
var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Perform operations on namespace",
	Long:  "Perform operations on namespace",
}

func init() {
	cmdv1.RootCmd.AddCommand(namespaceCmd)
}
