package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	resourceManagerLabel = "k8s-dev-training/resource-manager"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	managerName := os.Getenv("RESOURCE_MANAGER_NAME")
	if managerName == "" {
		return errors.New("RESOURCE_MANAGER_NAME is required")
	}

	commandJSON := os.Getenv("JOB_COMMAND")
	if commandJSON == "" {
		return errors.New("JOB_COMMAND is required")
	}

	var command []string
	if err := json.Unmarshal([]byte(commandJSON), &command); err != nil {
		return fmt.Errorf("parse JOB_COMMAND: %w", err)
	}
	if len(command) == 0 {
		return errors.New("JOB_COMMAND must not be empty")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("load in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create Kubernetes client: %w", err)
	}

	jobs := clientset.BatchV1().Jobs(namespace)
	jobName := managerName + "-job"

	// If this helper restarts, do not create another Job with a new name.
	_, err = jobs.Get(context.Background(), jobName, metav1.GetOptions{})
	if err == nil {
		log.Printf("Job %s/%s already exists", namespace, jobName)
		return waitForever()
	}
	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("check existing Job: %w", err)
	}

	backoffLimit := int32(1)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
			Labels: map[string]string{
				resourceManagerLabel: managerName,
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						resourceManagerLabel: managerName,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "task",
							Image:   "alpine",
							Command: command,
						},
					},
				},
			},
		},
	}

	_, err = jobs.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("create Job: %w", err)
	}

	log.Printf("Job %s/%s is present", namespace, jobName)
	// return waitForever()
	return nil
}

func waitForever() error {
	for {
		time.Sleep(time.Hour)
	}
}
