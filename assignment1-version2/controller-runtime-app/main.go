package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getKubeConfigRuntime() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func main() {
	ctx := context.Background()

	config, err := getKubeConfigRuntime()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Create unified controller-runtime client
	k8sClient, err := client.New(config, client.Options{})
	if err != nil {
		log.Fatalf("Failed to create controller-runtime client: %v", err)
	}

	ns := "default"

	// ==========================================
	// RESOURCE 1: ConfigMap CRUD
	// ==========================================
	fmt.Println("\n--- [controller-runtime] 1. ConfigMap CRUD ---")
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cr-demo-cm",
			Namespace: ns,
		},
		Data: map[string]string{"tier": "frontend"},
	}

	// CREATE
	if err := k8sClient.Create(ctx, cm); err != nil {
		log.Fatalf("Create failed: %v", err)
	}
	fmt.Println("Created ConfigMap: cr-demo-cm")

	// READ
	fetchedCM := &corev1.ConfigMap{}
	key := types.NamespacedName{Name: "cr-demo-cm", Namespace: ns}
	if err := k8sClient.Get(ctx, key, fetchedCM); err != nil {
		log.Fatalf("Get failed: %v", err)
	}
	fmt.Printf("Read ConfigMap tier: %s\n", fetchedCM.Data["tier"])

	// UPDATE
	fetchedCM.Data["tier"] = "backend"
	if err := k8sClient.Update(ctx, fetchedCM); err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	fmt.Println("Updated ConfigMap tier to backend")

	// DELETE
	if err := k8sClient.Delete(ctx, fetchedCM); err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	fmt.Println("Deleted ConfigMap: cr-demo-cm")

	// ==========================================
	// RESOURCE 2: Secret CRUD
	// ==========================================
	fmt.Println("\n--- [controller-runtime] 2. Secret CRUD ---")
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cr-demo-secret",
			Namespace: ns,
		},
		StringData: map[string]string{"api-key": "123456"},
	}

	// CREATE
	if err := k8sClient.Create(ctx, secret); err != nil {
		log.Fatalf("Create Secret failed: %v", err)
	}
	fmt.Println("Created Secret: cr-demo-secret")

	// READ
	fetchedSecret := &corev1.Secret{}
	secretKey := types.NamespacedName{Name: "cr-demo-secret", Namespace: ns}
	if err := k8sClient.Get(ctx, secretKey, fetchedSecret); err != nil {
		log.Fatalf("Get Secret failed: %v", err)
	}
	fmt.Printf("Read Secret key: %s\n", string(fetchedSecret.Data["api-key"]))

	// UPDATE
	fetchedSecret.StringData = map[string]string{"api-key": "654321"}
	if err := k8sClient.Update(ctx, fetchedSecret); err != nil {
		log.Fatalf("Update Secret failed: %v", err)
	}
	fmt.Println("Updated Secret key")

	// DELETE
	if err := k8sClient.Delete(ctx, fetchedSecret); err != nil {
		log.Fatalf("Delete Secret failed: %v", err)
	}
	fmt.Println("Deleted Secret: cr-demo-secret")

	// ==========================================
	// RESOURCE 3: Service CRUD
	// ==========================================
	fmt.Println("\n--- [controller-runtime] 3. Service CRUD ---")
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cr-demo-svc",
			Namespace: ns,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 80, Name: "http"},
			},
			Selector: map[string]string{"app": "demo"},
		},
	}

	// CREATE
	if err := k8sClient.Create(ctx, svc); err != nil {
		log.Fatalf("Create Service failed: %v", err)
	}
	fmt.Println("Created Service: cr-demo-svc")

	// READ
	fetchedSvc := &corev1.Service{}
	svcKey := types.NamespacedName{Name: "cr-demo-svc", Namespace: ns}
	if err := k8sClient.Get(ctx, svcKey, fetchedSvc); err != nil {
		log.Fatalf("Get Service failed: %v", err)
	}
	fmt.Printf("Read Service Port: %d\n", fetchedSvc.Spec.Ports[0].Port)

	// UPDATE
	fetchedSvc.Spec.Ports[0].Port = 8080
	if err := k8sClient.Update(ctx, fetchedSvc); err != nil {
		log.Fatalf("Update Service failed: %v", err)
	}
	fmt.Println("Updated Service Port to 8080")

	// DELETE
	if err := k8sClient.Delete(ctx, fetchedSvc); err != nil {
		log.Fatalf("Delete Service failed: %v", err)
	}
	fmt.Println("Deleted Service: cr-demo-svc")
}