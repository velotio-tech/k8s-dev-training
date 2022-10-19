package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var UseCtrlRuntime bool
var CtrlClient client.Client
var ClientSet *kubernetes.Clientset

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kube-client",
	Short: "Client for accessing kubernetes commands",
	Long:  "Client for accessing kubernetes commands",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if UseCtrlRuntime {
			fmt.Println("Using controller runtime")
			var err error
			CtrlClient, err = client.New(clientConfig.GetConfigOrDie(), client.Options{})
			if err != nil {
				fmt.Println("could not create client, error: ", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Using in-cluster config")
			config, err := rest.InClusterConfig()
			if err != nil {
				fmt.Println("Could not to create in-cluster config, trying to fetch from global kube config")
				kubeConfigFilepath := filepath.Join(
					os.Getenv("HOME"), ".kube", "config",
				)
				config, err = clientcmd.BuildConfigFromFlags("", kubeConfigFilepath)
				if err != nil {
					panic(err.Error())
				}
			}

			ClientSet, err = kubernetes.NewForConfig(config)
			if err != nil {
				fmt.Println("could not create clientset, error: ", err)
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Client for accessing kubernetes commands !!!\n")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&UseCtrlRuntime, "ctrl-runtime", "c", false, "Use controller runtime to create config")
}
