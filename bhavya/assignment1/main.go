package main

import (
	"fmt"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/api/configmap"
	deployment "github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/api/deployment"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/api/pod"
	"github.com/jnbhavya/k8s-dev-training/bhavya/assignment1/api/statefulset"
)

func main() {
	// client-go
	err := configmap.Operations()
	if err != nil {
		fmt.Println(err)
	}
	err =deployment.Operations()
	if err != nil {
		fmt.Println(err)
	}
	err = pod.Operations()
	if err != nil {
		fmt.Println(err)
	}
	err = statefulset.Operations()
	if err != nil {
		fmt.Println(err)
	}

	//runtime controller
	err = configmap.RtcOperation()
	if err != nil {
		fmt.Println(err)
	}
	err =deployment.RtcOperation()
	if err != nil {
		fmt.Println(err)
	}
	err = pod.RtcOperation()
	if err != nil {
		fmt.Println(err)
	}
	err = statefulset.RtcOperation()
	if err != nil {
		fmt.Println(err)
	}
}
