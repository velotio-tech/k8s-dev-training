package main

import (
	"context"
	"log"
	"sync"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// create the clientset
	clientset, err := createClientSet()
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()

	var wg sync.WaitGroup

	wg.Go(func() {
		podCRUDS(ctx, clientset)
	})

	wg.Go(func() {
		deploymentCRUDS(ctx, clientset)
	})

	wg.Go(func() {
		jobCRUDS(ctx, clientset)
	})

	wg.Wait()

	log.Println("All CRUD operations done")
}
