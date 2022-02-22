/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"assignment1/ops/c_go"
	"assignment1/ops/c_runtime"
	"log"

	"github.com/spf13/cobra"
)

// podCmd represents the pod command
var updateSvcCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if *name == "" {
			log.Fatal("Please enter a valid name for the resource to be updated.")
		}

		if *backend == "cgo" {
			err := c_go.UpdateService(*name, *namespace, *port)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := c_runtime.UpdateService(*name, *namespace, *port)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Print("Resource updated successfully.")
	},
}

func init() {
	updateCmd.AddCommand(updateSvcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// podCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// podCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}