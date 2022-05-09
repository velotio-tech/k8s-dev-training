package deployment

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDeployment(deploymentsClient v1.DeploymentInterface) string {
	// Create a new nginx deployment
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	depname := result.GetObjectMeta().GetName()
	fmt.Printf("Created deployment %q.\n", depname)
	return depname
}
func GetDeployment(deploymentsClient v1.DeploymentInterface, deploymentname string) appsv1.Deployment {
	// Return deployment spec, filter deployment by provided name
	result, getErr := deploymentsClient.Get(context.TODO(), deploymentname, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}
	return *result
}
func UpdateDeployment(deploymentsClient v1.DeploymentInterface, deploymentname string) {
	// Increase deployment replica to 2 and update labels
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), deploymentname, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}
		result.Spec.Template.Labels["application"] = "veloito"
		replicas := int32(2)
		result.Spec.Replicas = &replicas
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
}

func DeleteDeployment(deploymentsClient v1.DeploymentInterface, deploymentname string) {
	// Delete deployment specified by the deplpoyment name, don't wait till it gets completely removed
	delErr := deploymentsClient.Delete(context.TODO(), deploymentname, metav1.DeleteOptions{})
	if delErr != nil {
		panic(fmt.Errorf("Failed to delete latest version of Deployment: %v", delErr))
	}
}

func CreateService(controllerClient client.Client) string {
	// Create service port, not adding any selector
	// TODO(shubham): need more understanding of context and why its needed
	newService := apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-service",
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:     "port",
					Protocol: "TCP",
					Port:     80,
				},
			},
		},
	}
	err := controllerClient.Create(context.TODO(), &newService)
	if err != nil {
		panic(err)
	}
	return newService.ObjectMeta.Name
}

func ListServices(controllerClient client.Client) {
	// Print all services in the given namespace
	fmt.Println("Listing all services in the namespace")
	services := &apiv1.ServiceList{}
	err := controllerClient.List(context.Background(), services)
	if err != nil {
		panic(err)
	}
	fmt.Println(services)
}

func EditServiceByName(controllerClient client.Client, serviceName string) {
	fmt.Println("Editing service %s", serviceName)
	service := &apiv1.Service{}
	err := controllerClient.Get(context.TODO(), client.ObjectKey{
		Name: serviceName,
	}, service)
	if err != nil {
		panic(err)
	}
	service.Spec.Ports[0].Port = 443
	err = controllerClient.Update(context.TODO(), service)
	if err != nil {
		panic(err)
	}
}

func DeleteService(controllerClient client.Client, serviceName string) {
	// Delete service by name
	fmt.Println("Deleting service %s", serviceName)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
	}
	err := controllerClient.Delete(context.TODO(), service)
	if err != nil {
		panic(err)
	}
}
