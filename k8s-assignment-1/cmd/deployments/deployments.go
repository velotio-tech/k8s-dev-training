package deployments

import (
	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
)

// deploymentsCmd represents the deployments command
var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Perform operations on deployments",
	Long:  "Perform operations on deployments",
}

func init() {
	cmdv1.RootCmd.AddCommand(deploymentsCmd)
}
