package runtimectrl

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var namespace = "pankaj"
var podName = "demo-pod"

//PodOperations - perform CURD operations on pod object using controller runtime library
func PodOperations() {
	ctrClient := getClient()

	fileBytes, err := ioutil.ReadFile("configfile/pod.yaml")
	if err != nil {
		panic(err.Error())
	}

	var podSpec apiv1.Pod
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1024)

	if err := dec.Decode(&podSpec); err != nil {
		panic(err)
	}

	// Create Pod
	fmt.Println("Pod Create operation started ........")
	err = ctrClient.Create(context.Background(), &podSpec)
	if err != nil {
		panic(err)
	}
	time.Sleep(15 * time.Second)
	fmt.Println("Pod Create operation completed ........")

	// Get Pod
	fmt.Println("Pod Get operation started ........")
	err = ctrClient.Get(context.Background(), client.ObjectKey{Namespace: namespace,
		Name: podName,
	}, &podSpec)

	if err != nil {

		panic(fmt.Errorf("Failed to get latest version of Pod: %v", err.Error()))
	}

	fmt.Printf("Latest Pod: %s \n", podSpec.Name)
	fmt.Println("Pod Get operation completed ........")

	// Updating Pod
	fmt.Println("Pod Update operation started ........")
	podSpec.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	err = ctrClient.Update(context.Background(), &podSpec)
	if err != nil {

		panic(err.Error())
	}

	time.Sleep(15 * time.Second)
	fmt.Println("Pod Update operation completed ........")

	// Verifying Update
	fmt.Println("Pod Update verfication operation started ........")
	err = ctrClient.Get(context.Background(), client.ObjectKey{Namespace: namespace,
		Name: podName,
	}, &podSpec)

	if err != nil {
		panic(fmt.Errorf("Failed to get latest version of Pod: %v", err))
	}
	if podSpec.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Verified Successfully")
	} else {
		panic(fmt.Errorf("Verification failed. Image found %s, expected: nginx:1.13",
			podSpec.Spec.Containers[0].Image))
	}
	fmt.Println("Pod Update verfication operation completed ........")

	// Delete Pods
	fmt.Println("Pod Delete operation started ........")
	err = ctrClient.Delete(context.Background(), &podSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Pod Delete operation completed ........")
}
