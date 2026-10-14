package main

import (
	"context"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const defaultNamespace = "default"

func podCRUDS(ctx context.Context, clientset kubernetes.Interface) {
	podClient := clientset.CoreV1().Pods(defaultNamespace)

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crud-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "hello-world",
					Image: "hello-world:latest",
				},
			},
		},
	}

	// CREATE
	createdPod, err := podClient.Create(ctx, &pod, metav1.CreateOptions{})
	if err != nil {
		log.Println(err.Error())

		return
	}
	log.Println("Created Pod:", createdPod.Name)

	watcher, err := podClient.Watch(
		ctx,
		metav1.ListOptions{
			FieldSelector: "metadata.name=" + pod.Name,
		},
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer watcher.Stop()

	podRunningCtxWithTimeout, podRunningCtxWithTimeoutCancel := context.WithTimeout(ctx, 20*time.Second)
	defer podRunningCtxWithTimeoutCancel()
	err = WaitForCondition(
		podRunningCtxWithTimeout,
		watcher,
		func(pod *corev1.Pod) bool {
			return pod.Status.Phase == corev1.PodRunning
		},
	)

	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Pod is Running!")

	// GET
	getPod, err := podClient.Get(ctx, "crud-pod", metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Got Pod:", getPod.Name)

	// UPDATE
	getPod.Labels = map[string]string{
		"app": "crud-pod",
	}

	updatedPod, err := podClient.Update(ctx, getPod, metav1.UpdateOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Updated Pod:", updatedPod.Labels)

	// LIST
	pods, err := podClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Pod count:", len(pods.Items))

	// DELETE
	err = podClient.Delete(ctx, "crud-pod", metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}

	podDeleteCtxWithTimeout, podDeleteCtxWithTimeoutCancel := context.WithTimeout(ctx, 20*time.Second)
	defer podDeleteCtxWithTimeoutCancel()

	err = WaitForDeletion(podDeleteCtxWithTimeout, watcher)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Deleted Pod")
}

func deploymentCRUDS(ctx context.Context, clientset kubernetes.Interface) {
	deploymentClient := clientset.AppsV1().Deployments(defaultNamespace)

	replicas := int32(2)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crud-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "crud-deployment",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "crud-deployment",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}

	// CREATE
	createdDeployment, err := deploymentClient.Create(
		ctx,
		deployment,
		metav1.CreateOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Created Deployment:", createdDeployment.Name)

	watcher, err := deploymentClient.Watch(
		ctx,
		metav1.ListOptions{
			FieldSelector: "metadata.name=" + deployment.Name,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer watcher.Stop()

	deploymentRunningCtxWithTimeout, deploymentRunningCtxWithTimeoutCancel := context.WithTimeout(ctx, 20*time.Second)
	defer deploymentRunningCtxWithTimeoutCancel()
	err = WaitForCondition(
		deploymentRunningCtxWithTimeout,
		watcher,
		func(deployment *appsv1.Deployment) bool {
			return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
		},
	)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Deployment is ready!")
	// GET
	getDeployment, err := deploymentClient.Get(
		ctx,
		"crud-deployment",
		metav1.GetOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Got Deployment:", getDeployment.Name)

	// UPDATE
	*getDeployment.Spec.Replicas = 3

	updatedDeployment, err := deploymentClient.Update(
		ctx,
		getDeployment,
		metav1.UpdateOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Updated replicas:", *updatedDeployment.Spec.Replicas)

	// LIST
	deployments, err := deploymentClient.List(
		ctx,
		metav1.ListOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Deployment count:", len(deployments.Items))

	// DELETE
	deploymentDeleteCtxWithTimeout, deploymentDeleteCtxWithTimeoutCancel := context.WithTimeout(ctx, 20*time.Second)
	defer deploymentDeleteCtxWithTimeoutCancel()
	err = deploymentClient.Delete(
		deploymentDeleteCtxWithTimeout,
		"crud-deployment",
		metav1.DeleteOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Deleted Deployment")
}

func jobCRUDS(ctx context.Context, clientset kubernetes.Interface) {
	jobClient := clientset.BatchV1().Jobs(defaultNamespace)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crud-job",
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

	// CREATE
	createdJob, err := jobClient.Create(
		ctx,
		job,
		metav1.CreateOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Created Job:", createdJob.Name)

	watcher, err := jobClient.Watch(
		ctx,
		metav1.ListOptions{
			FieldSelector: "metadata.name=" + job.Name,
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	defer watcher.Stop()

	jobSuccessCtxWithTimeout, jobSuccessCtxWithTimeoutCancel := context.WithTimeout(ctx, 20*time.Second)
	defer jobSuccessCtxWithTimeoutCancel()
	err = WaitForCondition(
		jobSuccessCtxWithTimeout,
		watcher,
		func(job *batchv1.Job) bool {
			return job.Status.Succeeded > 0
		},
	)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Job completed!")
	// GET
	getJob, err := jobClient.Get(
		ctx,
		"crud-job",
		metav1.GetOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Got Job:", getJob.Name)

	// UPDATE
	getJob.Labels = map[string]string{
		"environment": "development",
	}

	updatedJob, err := jobClient.Update(
		ctx,
		getJob,
		metav1.UpdateOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Updated Job:", updatedJob.Labels)

	// LIST
	jobs, err := jobClient.List(
		ctx,
		metav1.ListOptions{},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Job count:", len(jobs.Items))

	// DELETE
	jobDeleteCtxWithTimeout, jobDeleteCtxWithTimeoutCancel := context.WithTimeout(ctx, 20*time.Second)
	defer jobDeleteCtxWithTimeoutCancel()
	err = jobClient.Delete(
		jobDeleteCtxWithTimeout,
		"crud-job",
		metav1.DeleteOptions{
			PropagationPolicy: new(metav1.DeletePropagationForeground),
		},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Deleted Job")
}
