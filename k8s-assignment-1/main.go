package main

import (
	"github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	_ "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd/deployments"
	_ "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd/namespace"
	_ "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd/pods"
)

func main() {
	cmd.Execute()
}
