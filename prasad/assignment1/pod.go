package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	podclient "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/client/pod/clientgo"
	podcrt "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/client/pod/crt"
	corev1 "k8s.io/api/core/v1"
)

var podOpsMutex sync.Mutex

func executePodOps() {
	blockingChannel := make(chan int)

	//	create pod
	go func() {
		var err error
		for {
			log.Println("creating pod...")
			podOpsMutex.Lock()
			err = podclient.CreatePod("my-pod", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create pod!!!")
			} else {
				log.Println("Pod created successfully...")
			}
			podOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	//	update pod
	go func() {
		var err error
		for {
			log.Println("Updating 'my-pod'")
			podOpsMutex.Lock()
			err = podclient.UpdatePod("my-pod", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create pod!!!")
			} else {
				log.Println("Pod updated successfully...")
			}
			podOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	//	Delete Pod
	go func() {
		var err error
		for {
			log.Println("Deleting pod...")
			podOpsMutex.Lock()
			err = podclient.DeletePod("my-pod", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to delete pod!!!")
			} else {
				log.Println("Pod deleted successfully...")
			}
			podOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	//	Print all pods
	go func() {
		for {
			log.Println("Printing all pods...")
			podOpsMutex.Lock()
			podList, err := podclient.GetAllPods(corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to fetch pods!!!")
			} else {
				for _, pod := range podList.Items {
					fmt.Println("Name - ", pod.GetName(), " Generate name - ", pod.GetGenerateName())
				}
			}
			podOpsMutex.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()

	<-blockingChannel
}

func executePodCrtOps() {
	blockingChannel := make(chan int)
	go func() {
		var err error
		for {
			log.Println("Creating pod...")
			podOpsMutex.Lock()
			err = podcrt.CreatePod("my-pod-crt", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create pod!!!")
			} else {
				log.Println("'my-pod-crt' created successfully...")
			}
			podOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Deleting pod...")
			podOpsMutex.Lock()
			err = podcrt.DeletePod("my-pod-crt", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to delete pod!!!")
			} else {
				log.Println("'my-pod-crt' deleted successfully...")
			}
			podOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Updating pod...")
			podOpsMutex.Lock()
			err = podcrt.UpdatePod("my-pod", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to update pod!!!")
			} else {
				log.Println("my-pod pod updated successfully...")
			}
			podOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			log.Println("Printing all pods...")
			podOpsMutex.Lock()
			podList, err := podcrt.GetAllPods(corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to fetch pods!!!")
			} else {
				for _, pod := range podList.Items {
					fmt.Println("Name - ", pod.GetName(), " Generate name - ", pod.GetGenerateName())
				}
			}
			podOpsMutex.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()

	<-blockingChannel
}
