package main

import (
	"context"
	"fmt"
	"time"

	"github.com/velotio-tech/k8s-dev-training/services"
	controllerruntime "github.com/velotio-tech/k8s-dev-training/services/controller-runtime"
	goclient "github.com/velotio-tech/k8s-dev-training/services/go-client"
)

func GetClient(name string) services.Client {
	switch name {
	case "go-client", "":
		return goclient.GetClient()
	case "controller-runtime":
		return controllerruntime.GetClient()
	}
	return nil
}

func main() {
	ctx := context.Background()
	client := GetClient("controller-runtime")
	ns := client.GetResource("ns")
	p := client.GetResource("pod")
	svc := client.GetResource("svc")

	// NAMESPACE OPERATIONS
	fmt.Println("=====CREATING NAMESPACE=====")
	err := ns.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====UPDATING NAMESPACE LABELS=====")
	err = ns.Update(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	err = ns.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// pod operations
	fmt.Println("=====CREATING POD=====")
	err = p.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====UPDATING POD LABELS=====")
	err = p.Update(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("")
	err = p.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// SERVICE OPERATIONS
	fmt.Println("=====CREATING SVC=====")
	err = svc.Create(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("=====UPDATING SVC SELECTORS=====")
	err = svc.Update(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("")
	err = svc.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// try to connect with check if nginx is reachable after port-forward, delays deletion
	time.Sleep(2 * time.Minute)

	fmt.Println("=====DELETING SVC=====")
	err = svc.Delete(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = svc.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("=====DELETING POD=====")
	err = p.Delete(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = p.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	// NAMESPACE OPERATIONS
	fmt.Println("=====DELETING NAMESPACE=====")
	err = ns.Delete(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	err = ns.List(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
