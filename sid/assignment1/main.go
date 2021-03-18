package main

import (
	"github.com/farkaskid/k8s-dev-training/assignment1/testers/deployment"
	"github.com/farkaskid/k8s-dev-training/assignment1/testers/pod"
	"github.com/farkaskid/k8s-dev-training/assignment1/testers/service"
)

func main() {
	// Go Client
	pod.GoClientTester()
	deployment.GoClientTester()
	service.GoClientTester()

	// Controller Runtime Client
	pod.CRTTester()
	deployment.CRTTester()
	service.CRTTester()
}
