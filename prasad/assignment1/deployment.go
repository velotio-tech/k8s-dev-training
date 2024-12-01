package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	depclient "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/client/deployment/clientgo"
	depcrtclient "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/client/deployment/crt"
	corev1 "k8s.io/api/core/v1"
)

var depOpsMutex sync.Mutex

func executeDeploymentOps() {
	go func() {
		var err error
		for {
			log.Println("Creating deployment...")
			depOpsMutex.Lock()
			err = depclient.CreateDeployment("inc-deployment", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create deployement!!!")
			} else {
				log.Println("'inc-deployment' created successfully...")
			}
			depOpsMutex.Unlock()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Deleting deployment...")
			depOpsMutex.Lock()
			err = depclient.DeleteDeployment("inc-deployment", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to delete deployment!!!")
			} else {
				log.Println("'inc-deployment' deleted successfully...")
			}
			depOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Updating deployment...")
			depOpsMutex.Lock()
			err = depclient.UpdateDeployment("inc-deployment", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to update deployment!!!")
			} else {
				log.Println("'inc-deployment' updated successfully...")
			}
			depOpsMutex.Unlock()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		var err error
		log.Println("All deployments...")
		depOpsMutex.Lock()
		deployments, err := depclient.GetAllDeployments(corev1.NamespaceDefault)
		if err != nil {
			log.Println(err, ":: Unable to fetch deployments!!!")
		}
		depOpsMutex.Unlock()
		for _, dep := range deployments.Items {
			fmt.Println("Name - ", dep.GetName())
		}
	}()
}

func executeDeploymentCrtOps() {
	go func() {
		var err error
		for {
			log.Println("Creating deployment...")
			depOpsMutex.Lock()
			err = depcrtclient.CreateDeployment("crt-deployment", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create deployment!!!")
			} else {
				log.Println("crt-deployment created successfully...")
			}
			depOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Deleting crt deployment...")
			depOpsMutex.Lock()
			err = depcrtclient.DeleteDeployment("crt-deployment", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to delete deployment!!!")
			} else {
				log.Println("crt-deployment deleted successfully...")
			}
			depOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Updating deployment...")
			depOpsMutex.Lock()
			err = depcrtclient.UpdateDeployment("crt-deployment", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to update deployment!!!")
			} else {
				log.Println("crt-deployment updated successfully...")
			}
			depOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			depOpsMutex.Lock()
			log.Println("Deployment List...")
			depList, err := depcrtclient.GetAllDeployments(corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to fetch all deployments!!!")
			} else {
				log.Println("All Deployments...")
				for _, dep := range depList.Items {
					fmt.Println("Name - ", dep.Name, " Generate-Name - ", dep.GetGenerateName())
				}
			}
			depOpsMutex.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()
}
