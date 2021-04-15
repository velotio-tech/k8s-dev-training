package main
import (
	"fmt"
	"context"

	// appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helpers "github.com/swapnil-velotio/k8s-dev-training/swapnil/helpers"

)

func main() {
	clientset := helpers.GetClientset()
	configMapCli := clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault)
	configMap := &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-configmap",
		},
		Data: map[string]string{"NAME":"Swapnil", "MOBILE": "8888888888"},
	}
	// creating pod
	fmt.Println("Enter to create configmap")
	helpers.Prompt()
	fmt.Println("Creating configmap...")
	result, err := configMapCli.Create(context.TODO(), configMap, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created map %q.\n", result.GetObjectMeta().GetName())

	fmt.Println("Enter to update cm")
	helpers.Prompt()
	fmt.Println("updating cm")
	cm, getErr := configMapCli.Get(context.TODO(), "demo-configmap", metav1.GetOptions{})
	if getErr != nil {
		panic("couldn't get pod details")
	}
	cm.Data["MOBILE"] = "9999999999"
	_, updateErr := configMapCli.Update(context.TODO(), cm, metav1.UpdateOptions{})
	if updateErr != nil {
		panic(updateErr)
	}
	fmt.Println("updated cm...")

	fmt.Println("Enter to get cms")
	helpers.Prompt()
	fmt.Println("getting cms")
	pods, listPodsErr := configMapCli.List(context.TODO(), metav1.ListOptions{})
	if listPodsErr != nil {
		panic(listPodsErr)
	}
	for _, d := range pods.Items {
		fmt.Printf(" * cm name - %s \n", d.Name)
	}
	fmt.Printf("cms fetched\n")

	fmt.Println("Enter to delete cm")
	helpers.Prompt()
	fmt.Println("deleting cm")
	deletePolicy := metav1.DeletePropagationForeground
	if err := configMapCli.Delete(context.TODO(), "demo-configmap", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Printf("cm deleted\n")





}