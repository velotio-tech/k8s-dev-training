package deployments

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

var deployment = &appsv1.Deployment {
	ObjectMeta: metav1.ObjectMeta {
		Name: "my-deployment",
	},
	Spec: appsv1.DeploymentSpec {
		Replicas: int32Ptr(2),
		Selector: &metav1.LabelSelector {
			MatchLabels: map[string]string {
				"app": "demo",
			},
		},
		Template: apiv1.PodTemplateSpec {
			ObjectMeta: metav1.ObjectMeta {
				Labels: map[string]string {
					"app": "demo",
				},
			},
			Spec: apiv1.PodSpec {
				Containers: []apiv1.Container {
					{
						Name:  "web",
						Image: "nginx",
						Ports: []apiv1.ContainerPort {
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

var deploymentClient v1.DeploymentInterface

func int32Ptr(i int32) *int32 { return &i }

func CreateDeploymentClient(clientset *kubernetes.Clientset)  {
	deploymentClient = clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
}

func GetAllDeployments() {
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func CreateDeployment() {
	fmt.Println("Creating deployment...", deploymentClient, deployment)
	result, err := deploymentClient.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Println("Error Occurred while creating the deployment", err.Error())
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func DeleteDeployment(){
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentClient.Delete(context.Background(), "my-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

func UpdateDeployment() {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentClient.Get(context.Background(), "my-deployment", metav1.GetOptions{})
		if getErr != nil {
			log.Println(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(1)                           // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentClient.Update(context.Background(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
}