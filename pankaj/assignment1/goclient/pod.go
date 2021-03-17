package goclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var namespace = "pankaj"
var podName = "demo-pod"

//PodOperations - perform CURD operations on pod object using client-go library
func PodOperations() {
	clientset := KubeConfig()
	podsClient := clientset.CoreV1().Pods(namespace)

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
	pods, err := podsClient.Create(context.TODO(), &podSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	time.Sleep(15 * time.Second) //wait to create pod
	fmt.Printf("Created pod %q.\n", pods.GetObjectMeta().GetName())
	fmt.Println("Pod Create operation completed ........")

	// Get Pod
	fmt.Println("Pod Get operation started ........")
	pods, getErr := podsClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if getErr != nil {

		panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
	}
	fmt.Printf("Latest Pod: %s \n", pods.Name)
	fmt.Println("Pod Get operation completed ........")

	// Updating Pod
	fmt.Println("Pod Update operation started ........")
	pods.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
	_, updateErr := podsClient.Update(context.TODO(), pods, metav1.UpdateOptions{})
	if updateErr != nil {
		panic(updateErr)
	}
	fmt.Println("Pod Update operation completed ........")

	time.Sleep(30 * time.Second)

	// Verifying Update
	fmt.Println("Pod Update verification operation started ........")
	pods, getErr = podsClient.Get(context.TODO(), podName, metav1.GetOptions{})
	if getErr != nil {

		panic(fmt.Errorf("Failed to get latest version of Pod: %v", getErr))
	}
	if pods.Spec.Containers[0].Image == "nginx:1.13" {
		fmt.Println("Verified Successfully")
	} else {

		panic(fmt.Errorf("Verification failed. Image found %s, expected: nginx:1.13",
			pods.Spec.Containers[0].Image))
	}
	fmt.Println("Pod update verfication operation completed ........")

	// Delete Pods
	fmt.Println("Pod Delete operation started ........")
	deletePolicy := metav1.DeletePropagationForeground
	err = podsClient.Delete(context.TODO(), podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pod Delete operation completed ........")
}
