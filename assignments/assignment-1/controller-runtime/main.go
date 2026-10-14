package main

import (
	"context"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	scheme := runtime.NewScheme()

	if err := batchv1.AddToScheme(scheme); err != nil {
		log.Fatal(err)
	}

	c, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatal(err)
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "controller-runtime-job",
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "crud-job",
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  "hello-world",
							Image: "hello-world:latest",
						},
					},
				},
			},
		},
	}

	ctx := context.Background()

	if err := c.Create(ctx, job); err != nil {
		log.Fatal(err)
	}

	if err = c.Get(ctx, client.ObjectKeyFromObject(job), job); err != nil {
		log.Fatal(err)
	}

	log.Println("Get job name:", job.GetName())

	requirement, err := labels.NewRequirement(
		"app",
		selection.Equals,
		[]string{"crud-job"},
	)
	if err != nil {
		log.Fatal(err)
	}

	lablesSelector := labels.NewSelector()
	lablesSelector.Add(*requirement)

	var jobs batchv1.JobList
	if err = c.List(ctx, &jobs, &client.ListOptions{
		Namespace:     job.GetNamespace(),
		LabelSelector: lablesSelector,
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("Number of jobs:", len(jobs.Items))

	jobLabels := job.GetLabels()
	jobLabels["env"] = "dev"
	job.SetLabels(jobLabels)

	if err = c.Update(ctx, job); err != nil {
		log.Fatal(err)
	}

	log.Println("updated job name:", job.GetName())

	if err = c.Delete(ctx, job); err != nil {
		log.Fatal(err)
	}
	log.Println("Job Deleted")
}
