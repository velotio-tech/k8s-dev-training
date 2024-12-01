package deployments

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

var deployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: "demo-deployment",
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(2),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "demo",
			},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": "demo",
				},
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  "web",
						Image: "nginx:1.12",
						Ports: []apiv1.ContainerPort{
							{
								Name:          "http",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 80,
							},
						},
					},
				},
			},
		},
	},
}
var deploymentsClient v1.DeploymentInterface

func CreateDeploymentClient(clientset *kubernetes.Clientset) {
	deploymentsClient = clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
}

func int32Ptr(i int32) *int32 { return &i }

func CreateDeployment() {
	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	// if err != nil {
	// 	panic(err)
	// }
	if err != nil {
		log.Println("Error occcured while creating deploymnent", err.Error())
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}
func ListDeployment() {
	// List Deployments
	//prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}
func DeleteDeployment() {
	// Delete Deployment
	//prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.Background(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

func UpdateDeployment() {
	// Update Deployment
	//prompt()
	fmt.Println("Updating deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(context.Background(), "demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			//panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
			log.Println(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentsClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
}

// func prompt() {
// 	fmt.Printf("-> Press Return key to continue.")
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		break
// 	}
// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}
// 	fmt.Println()
// }
