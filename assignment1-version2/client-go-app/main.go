package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// getKubeConfig returns in-cluster config if running in a Pod, or falls back to kubeconfig.
func getKubeConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err == nil {
		fmt.Println("Using In-Cluster Config")
		return config, nil
	}

	fmt.Println("Falling back to local Kubeconfig")
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func main() {
	ctx := context.Background()

	config, err := getKubeConfig()
	if err != nil {
		log.Fatalf("Failed to build config: %v", err)
	}

	// Create typed clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	ns := "default"

	// ==========================================
	// RESOURCE 1: ConfigMap CRUD
	// ==========================================
	fmt.Println("\n--- [client-go] 1. ConfigMap CRUD ---")
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-configmap"},
		Data:       map[string]string{"env": "dev"},
	}

	// CREATE
	createdCM, err := clientset.CoreV1().ConfigMaps(ns).Create(ctx, cm, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Create CM failed: %v", err)
	}
	fmt.Printf("Created ConfigMap: %s\n", createdCM.Name)
	time.Sleep(2 * time.Second) // Wait for the resource to be created

	// READ
	getCM, err := clientset.CoreV1().ConfigMaps(ns).Get(ctx, "demo-configmap", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Get CM failed: %v", err)
	}
	fmt.Printf("Read ConfigMap env value: %s\n", getCM.Data["env"])
	time.Sleep(2 * time.Second) // Wait for the resource to be read

	// UPDATE
	getCM.Data["env"] = "production"
	updatedCM, err := clientset.CoreV1().ConfigMaps(ns).Update(ctx, getCM, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Update CM failed: %v", err)
	}
	fmt.Printf("Updated ConfigMap env value: %s\n", updatedCM.Data["env"])
	time.Sleep(2 * time.Second) // Wait for the resource to be updated

	// DELETE
	err = clientset.CoreV1().ConfigMaps(ns).Delete(ctx, "demo-configmap", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Delete CM failed: %v", err)
	}
	fmt.Println("Deleted ConfigMap: demo-configmap")
	time.Sleep(2 * time.Second) // Wait for the resource to be deleted

	// ==========================================
	// RESOURCE 2: Secret CRUD
	// ==========================================
	fmt.Println("\n--- [client-go] 2. Secret CRUD ---")
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-secret"},
		StringData: map[string]string{"password": "super-secret-pass"},
	}

	// CREATE
	_, err = clientset.CoreV1().Secrets(ns).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Create Secret failed: %v", err)
	}
	fmt.Println("Created Secret: demo-secret")
	time.Sleep(4 * time.Second) // Wait for the resource to be created

	// READ
	getSecret, err := clientset.CoreV1().Secrets(ns).Get(ctx, "demo-secret", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Get Secret failed: %v", err)
	}
	fmt.Printf("Read Secret password: %s\n", string(getSecret.Data["password"]))
	time.Sleep(4 * time.Second) // Wait for the resource to be read

	// UPDATE
	getSecret.StringData = map[string]string{"password": "new-secret-pass"}
	_, err = clientset.CoreV1().Secrets(ns).Update(ctx, getSecret, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Update Secret failed: %v", err)
	}
	fmt.Println("Updated Secret password")
	time.Sleep(4 * time.Second) // Wait for the resource to be updated

	// DELETE
	err = clientset.CoreV1().Secrets(ns).Delete(ctx, "demo-secret", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Delete Secret failed: %v", err)
	}
	fmt.Println("Deleted Secret: demo-secret")
	time.Sleep(4 * time.Second) // Wait for the resource to be deleted

	// ==========================================
	// RESOURCE 3: Namespace CRUD
	// ==========================================
	fmt.Println("\n--- [client-go] 3. Namespace CRUD ---")
	testNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "demo-test-ns",
			Labels: map[string]string{"team": "devops"},
		},
	}

	// CREATE
	_, err = clientset.CoreV1().Namespaces().Create(ctx, testNS, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Create Namespace failed: %v", err)
	}
	fmt.Println("Created Namespace: demo-test-ns")

	// READ
	getNS, err := clientset.CoreV1().Namespaces().Get(ctx, "demo-test-ns", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Get Namespace failed: %v", err)
	}
	fmt.Printf("Read Namespace label 'team': %s\n", getNS.Labels["team"])

	// UPDATE
	getNS.Labels["team"] = "platform"
	_, err = clientset.CoreV1().Namespaces().Update(ctx, getNS, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Update Namespace failed: %v", err)
	}
	fmt.Println("Updated Namespace label 'team' to platform")

	// DELETE
	err = clientset.CoreV1().Namespaces().Delete(ctx, "demo-test-ns", metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Delete Namespace failed: %v", err)
	}
	fmt.Println("Deleted Namespace: demo-test-ns")
}