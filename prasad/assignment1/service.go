package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	svcclient "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/client/service/clientgo"
	svc_crt "github.com/thisisprasad/k8s-dev-training/prasad/assignment1/client/service/crt"

	corev1 "k8s.io/api/core/v1"
)

var svcOpsMutex sync.Mutex

func executeSvcOps() {
	go func() {
		var err error
		for {
			log.Println("Creating service...")
			svcOpsMutex.Lock()
			err = svcclient.CreateService("my-svc", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create service...")
			} else {
				log.Println("my-svc created successfully...")
			}
			svcOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Deleting service...")
			svcOpsMutex.Lock()
			err = svcclient.DeleteService("my-svc", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to delete svc!!!")
			} else {
				log.Println("my-svc deleted...")
			}
			svcOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Updating service...")
			svcOpsMutex.Lock()
			err = svcclient.UpdateService("my-svc", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to update service!!!")
			} else {
				log.Println("my-svc updated successfully...")
			}
			svcOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			log.Println("Printing all services...")
			svcOpsMutex.Lock()
			svcList, err := svcclient.GetAllServices(corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to fetch services!!!")
			} else {
				for _, svc := range svcList.Items {
					fmt.Println("Name - ", svc.GetName(), " Generate name - ", svc.GetGenerateName())
				}
			}
			svcOpsMutex.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()
}

func executeSvcCrtOps() {
	go func() {
		var err error
		for {
			log.Println("Creating service...")
			svcOpsMutex.Lock()
			err = svc_crt.CreateService("my-svc-crt", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to create service!!!")
			} else {
				log.Println("my-svc-crt service created successfully...")
			}
			svcOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Deleting service...")
			svcOpsMutex.Lock()
			err = svc_crt.DeleteService("my-svc-crt", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to delete service")
			} else {
				log.Println("my-svc-crt service deleted successfully...")
			}
			svcOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		var err error
		for {
			log.Println("Updating service...")
			svcOpsMutex.Lock()
			err = svc_crt.UpdateService("my-svc-crt", corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to update service!!!")
			} else {
				log.Println("my-svc-crt service updated successfully...")
			}
			svcOpsMutex.Unlock()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			log.Println("Print all services...")
			svcOpsMutex.Lock()
			svcList, err := svc_crt.GetAllServices(corev1.NamespaceDefault)
			if err != nil {
				log.Println(err, ":: Unable to fetch services!!!")
			} else {
				for _, svc := range svcList.Items {
					fmt.Println("Name - ", svc.GetName(), " Generate name - ", svc.GetGenerateName())
				}
			}
			svcOpsMutex.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()
}
