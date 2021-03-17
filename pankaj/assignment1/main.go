package main

import (
	"fmt"

	"github.com/pankaj9310/k8s-dev-training/pankaj/assignment1/goclient"
	"github.com/pankaj9310/k8s-dev-training/pankaj/assignment1/runtimectrl"
)

func main() {

	fmt.Println("GO client operation start")

	goclient.PodOperations()
	goclient.DeploymentOperations()
	goclient.RoleOperations()
	fmt.Println("GO client operation completed")

	fmt.Printf("\n ******************* \n")

	fmt.Println("Runtime Controller operation start")
	runtimectrl.PodOperations()
	runtimectrl.DeploymentOperations()
	runtimectrl.RoleOperations()
	fmt.Println("Runtime Controller operation completed")
}
