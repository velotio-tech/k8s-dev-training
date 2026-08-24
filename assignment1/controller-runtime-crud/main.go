package main

// Same three CRUD exercises as the client-go program, but using
// controller-runtime's unified client.Client interface. This is the same
// interface you will use inside every Reconciler you write for the rest
// of the assignments, so getting comfortable with it now pays off later.
//
// Notice there is no ConfigMapClient / PodClient / DeploymentClient - it's
// one Get/List/Create/Update/Delete for everything, and the object type
// passed in decides which resource is affected.

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// ctrl.GetConfig() tries in-cluster config first, then falls back to
	// your local kubeconfig (~/.kube/config or $KUBECONFIG) automatically -
	// controller-runtime does this fallback for you, unlike client-go
	// where we wrote it by hand in the other program.
	cfg, err := ctrl.GetConfig()
	check(err)

	// controller-runtime needs a Scheme so it knows which Go struct maps
	// to which Kubernetes Kind. scheme.Scheme (from client-go) already
	// has ConfigMap, Pod, Deployment etc registered - this is the same
	// registration mechanism you will extend in Assignment 2 to add your
	// own CRD type.
	s := scheme.Scheme

	// client.New gives a client that talks directly to the API server
	// (no cache) - good for a one-shot program like this. Inside a real
	// controller (built with ctrl.NewManager), the manager gives you a
	// client that reads from a cache instead - see the README note.
	c, err := client.New(cfg, client.Options{Scheme: s})
	check(err)

	ctx := context.Background()
	ns := "default"

	// ---------------------------------------------------------------
	// ConfigMap CRUD
	// ---------------------------------------------------------------
	fmt.Println("\n== ConfigMap CRUD (controller-runtime) ==")

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-cm-cr", Namespace: ns},
		Data:       map[string]string{"key1": "value1"},
	}
	check(c.Create(ctx, cm))
	fmt.Println("created demo-cm-cr")

	gotCM := &corev1.ConfigMap{}
	check(c.Get(ctx, types.NamespacedName{Name: "demo-cm-cr", Namespace: ns}, gotCM))
	fmt.Println("read demo-cm-cr data:", gotCM.Data)

	gotCM.Data["key2"] = "value2"
	check(c.Update(ctx, gotCM))
	fmt.Println("updated demo-cm-cr data:", gotCM.Data)

	check(c.Delete(ctx, gotCM))
	fmt.Println("deleted demo-cm-cr")

	// ---------------------------------------------------------------
	// Pod CRUD
	// ---------------------------------------------------------------
	fmt.Println("\n== Pod CRUD (controller-runtime) ==")

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-pod-cr", Namespace: ns},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "nginx", Image: "nginx:alpine"},
			},
		},
	}
	check(c.Create(ctx, pod))
	fmt.Println("created demo-pod-cr")

	gotPod := &corev1.Pod{}
	check(c.Get(ctx, types.NamespacedName{Name: "demo-pod-cr", Namespace: ns}, gotPod))
	fmt.Println("read demo-pod-cr phase:", gotPod.Status.Phase)

	if gotPod.Labels == nil {
		gotPod.Labels = map[string]string{}
	}
	gotPod.Labels["env"] = "test"
	check(c.Update(ctx, gotPod))
	fmt.Println("updated demo-pod-cr with label env=test")

	check(c.Delete(ctx, gotPod))
	fmt.Println("deleted demo-pod-cr")

	// ---------------------------------------------------------------
	// Deployment CRUD
	// ---------------------------------------------------------------
	fmt.Println("\n== Deployment CRUD (controller-runtime) ==")

	replicas := int32(1)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-dep-cr", Namespace: ns},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "demo-cr"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "demo-cr"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "nginx", Image: "nginx:alpine"},
					},
				},
			},
		},
	}
	check(c.Create(ctx, dep))
	fmt.Println("created demo-dep-cr")

	gotDep := &appsv1.Deployment{}
	check(c.Get(ctx, types.NamespacedName{Name: "demo-dep-cr", Namespace: ns}, gotDep))
	fmt.Println("read demo-dep-cr replicas:", *gotDep.Spec.Replicas)

	newReplicas := int32(2)
	gotDep.Spec.Replicas = &newReplicas
	check(c.Update(ctx, gotDep))
	fmt.Println("scaled demo-dep-cr to 2 replicas")

	check(c.Delete(ctx, gotDep))
	fmt.Println("deleted demo-dep-cr")

	fmt.Println("\nAll controller-runtime CRUD operations completed successfully")
}
