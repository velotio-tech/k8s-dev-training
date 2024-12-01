package main

import "log"

func main() {
	waitChan := make(chan int)
	log.Println("starting k8s resources ops...")
	go executePodOps()
	go executePodCrtOps()
	go executeDeploymentOps()
	go executeDeploymentCrtOps()
	go executeSvcOps()
	go executeSvcCrtOps()

	<-waitChan
}
