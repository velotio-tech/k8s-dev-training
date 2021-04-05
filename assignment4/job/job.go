package main

import (
	"context"
	"fmt"
	"time"

	"os"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	newclient, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("Failed to create client")
		os.Exit(1)
	}

	const NamespaceName string = "pankaj"

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pankaj-job-" + fmt.Sprintf("%v", time.Now().Unix()),
			Namespace: NamespaceName,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pankaj-job",
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "pankaj-job",
							Image: "busybox",
							Command: []string{
								"/bin/sh",
								"-ec",
								"sleep 5",
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyOnFailure,
				},
			},
		},
	}

	err = newclient.Create(context.TODO(), job)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job Created successfully!")

	fmt.Println("Listing job: ")
	jobList := batchv1.JobList{}
	err = newclient.List(context.TODO(), &jobList)
	if err != nil {
		panic(err.Error())
	}
	for _, jobs := range jobList.Items {
		fmt.Println("This is the job: ", jobs.Name)
	}

	// fmt.Println("Deleting Job ... ")
	// if deleteErr := newclient.Delete(context.TODO(), job); deleteErr != nil {
	// 	panic(err.Error())
	// }
	// fmt.Println("Job deleted")
}
