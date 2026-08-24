package main

// This program uses the raw client-go library to perform CRUD operations
// on three resource types: ConfigMap, Pod, and Deployment.
//
// client-go gives you one "typed client" per resource, e.g.
// clientset.CoreV1().Pods(ns) - every call is a direct REST call to the
// API server. There is no caching unless you build an informer yourself.

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// getClientset builds a Kubernetes clientset.
// The assignment specifically asks for "in-cluster config", which is what
// you get automatically when this binary runs as a Pod inside the cluster
// (via the ServiceAccount token mounted at /var/run/secrets/...).
//
// Since we are developing locally against a kind cluster (not yet deployed
// as a Pod), we fall back to the local kubeconfig so you can iterate fast.
// Once this works, see the README for how to actually run it in-cluster.
func getClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		fmt.Println("[info] using in-cluster config")
		return kubernetes.NewForConfig(config)
	}

	fmt.Println("[info] not running in-cluster, falling back to local kubeconfig")
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	flag.Parse()
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("could not load kubeconfig: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	clientset, err := getClientset()
	check(err)

	ctx := context.Background()
	ns := "default"

	// ---------------------------------------------------------------
	// ConfigMap CRUD
	// ---------------------------------------------------------------
	fmt.Println("\n== ConfigMap CRUD (client-go) ==")

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-cm"},
		Data:       map[string]string{"key1": "value1"},
	}

	// CREATE
	_, err = clientset.CoreV1().ConfigMaps(ns).Create(ctx, cm, metav1.CreateOptions{})
	check(err)
	fmt.Println("created demo-cm")

	// READ
	gotCM, err := clientset.CoreV1().ConfigMaps(ns).Get(ctx, "demo-cm", metav1.GetOptions{})
	check(err)
	fmt.Println("read demo-cm data:", gotCM.Data)

	// UPDATE (add a second key)
	gotCM.Data["key2"] = "value2"
	updatedCM, err := clientset.CoreV1().ConfigMaps(ns).Update(ctx, gotCM, metav1.UpdateOptions{})
	check(err)
	fmt.Println("updated demo-cm data:", updatedCM.Data)

	// DELETE
	err = clientset.CoreV1().ConfigMaps(ns).Delete(ctx, "demo-cm", metav1.DeleteOptions{})
	check(err)
	fmt.Println("deleted demo-cm")

	// ---------------------------------------------------------------
	// Pod CRUD
	// ---------------------------------------------------------------
	fmt.Println("\n== Pod CRUD (client-go) ==")

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "nginx", Image: "nginx:alpine"},
			},
		},
	}

	_, err = clientset.CoreV1().Pods(ns).Create(ctx, pod, metav1.CreateOptions{})
	check(err)
	fmt.Println("created demo-pod")

	gotPod, err := clientset.CoreV1().Pods(ns).Get(ctx, "demo-pod", metav1.GetOptions{})
	check(err)
	fmt.Println("read demo-pod phase:", gotPod.Status.Phase)

	// UPDATE - pods only allow a limited set of mutable fields
	// (labels, some annotations, tolerations grow-only, image for
	// certain fields via ephemeral containers). We update a label here.
	if gotPod.Labels == nil {
		gotPod.Labels = map[string]string{}
	}
	gotPod.Labels["env"] = "test"
	_, err = clientset.CoreV1().Pods(ns).Update(ctx, gotPod, metav1.UpdateOptions{})
	check(err)
	fmt.Println("updated demo-pod with label env=test")

	err = clientset.CoreV1().Pods(ns).Delete(ctx, "demo-pod", metav1.DeleteOptions{})
	check(err)
	fmt.Println("deleted demo-pod")

	// ---------------------------------------------------------------
	// Deployment CRUD
	// ---------------------------------------------------------------
	fmt.Println("\n== Deployment CRUD (client-go) ==")

	replicas := int32(1)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-dep"},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "demo"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "demo"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:alpine"},
					},
				},
			},
		},
	}

	_, err = clientset.AppsV1().Deployments(ns).Create(ctx, dep, metav1.CreateOptions{})
	check(err)
	fmt.Println("created demo-dep")

	gotDep, err := clientset.AppsV1().Deployments(ns).Get(ctx, "demo-dep", metav1.GetOptions{})
	check(err)
	fmt.Println("read demo-dep replicas:", *gotDep.Spec.Replicas)

	// UPDATE - scale to 2 replicas
	newReplicas := int32(2)
	gotDep.Spec.Replicas = &newReplicas
	_, err = clientset.AppsV1().Deployments(ns).Update(ctx, gotDep, metav1.UpdateOptions{})
	check(err)
	fmt.Println("scaled demo-dep to 2 replicas")

	err = clientset.AppsV1().Deployments(ns).Delete(ctx, "demo-dep", metav1.DeleteOptions{})
	check(err)
	fmt.Println("deleted demo-dep")

	fmt.Println("\nAll client-go CRUD operations completed successfully")
}
