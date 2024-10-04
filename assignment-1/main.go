package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func int32Ptr(i int32) *int32 { return &i }

type Resource struct {
	name  string
	ready bool
}
type K8sClient interface {
	GetResource() (*Resource, error)
	CreateResource() error
	UpdateResource() error
	DeleteResource() error
}

type DeploymentClientSet struct {
	clientSet  *kubernetes.Clientset
	deployment *appsv1.Deployment
}

func IsDeploymentReady(deployment *appsv1.Deployment) bool {
	availableReplicas := int32(0)
	// Check if the deployment is progressing
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == apiv1.ConditionTrue {
			fmt.Println(condition.Message)
		}
		if condition.Type == appsv1.DeploymentAvailable && condition.Status == apiv1.ConditionTrue {
			fmt.Println(condition.Message)
			availableReplicas += 1
		}
	}
	return availableReplicas == deployment.Status.Replicas
}

func ToDeploymentResource(deployment *appsv1.Deployment) (*Resource, error) {
	return &Resource{
		name:  deployment.Name,
		ready: IsDeploymentReady(deployment),
	}, nil
}

func (d *DeploymentClientSet) CreateResource() error {
	result, err := d.clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Create(context.TODO(), d.deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deployment %s created\n", result.Name)
	return nil
}

func (d *DeploymentClientSet) GetResource() (*Resource, error) {
	result, err := d.clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Get(context.TODO(), d.deployment.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	d.deployment = result
	return ToDeploymentResource(result)
}

func (d *DeploymentClientSet) UpdateResource() error {
	fmt.Println("Trigger Deployment update")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := d.clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Get(context.TODO(), d.deployment.Name, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		result.Spec.Replicas = int32Ptr(1)
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
		_, updateErr := d.clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}
	return nil
}

func (d *DeploymentClientSet) DeleteResource() error {
	deletePolicy := metav1.DeletePropagationForeground
	err := d.clientSet.AppsV1().Deployments(apiv1.NamespaceDefault).Delete(context.TODO(), d.deployment.Name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	fmt.Println("Deployment deleted")
	return nil
}

type ServiceClientSet struct {
	clientSet *kubernetes.Clientset
	service   *apiv1.Service
}

func IsServiceReady(service *apiv1.Service) bool {
	ready := true
	// Access the Service ClusterIP
	fmt.Printf("Service %s found with ClusterIP: %s\n", service.Name, service.Spec.ClusterIP)

	// If it’s a LoadBalancer type service, check external IPs
	if service.Spec.Type == apiv1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			fmt.Printf("Service %s is exposed via external IP: %s\n", service.Name, service.Status.LoadBalancer.Ingress[0].IP)
		} else {
			fmt.Printf("Service %s is waiting for an external IP\n", service.Name)
			ready = false
		}
	}
	return ready
}

func ToServiceResource(service *apiv1.Service) (*Resource, error) {
	return &Resource{
		name:  service.Name,
		ready: IsServiceReady(service),
	}, nil
}

func (s *ServiceClientSet) CreateResource() error {
	result, err := s.clientSet.CoreV1().Services(apiv1.NamespaceDefault).Create(context.TODO(), s.service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Service %s created\n", result.Name)
	return nil
}

func (s *ServiceClientSet) GetResource() (*Resource, error) {
	result, err := s.clientSet.CoreV1().Services(apiv1.NamespaceDefault).Get(context.TODO(), s.service.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	s.service = result
	return ToServiceResource(result)
}

func (s *ServiceClientSet) UpdateResource() error {
	fmt.Println("Service update is not implemented")
	return nil
}

func (s *ServiceClientSet) DeleteResource() error {
	err := s.clientSet.CoreV1().Services(apiv1.NamespaceDefault).Delete(context.TODO(), s.service.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Service deleted")
	return nil
}

type VolumeClientSet struct {
	clientSet *kubernetes.Clientset
	pvc       *apiv1.PersistentVolumeClaim
}

func IsVolumeReady(pvc *apiv1.PersistentVolumeClaim) bool {
	ready := false
	// Access the PVC status
	fmt.Printf("PVC %s found with status phase: %s\n", pvc.Name, pvc.Status.Phase)

	// Check if PVC is bound
	if pvc.Status.Phase == apiv1.ClaimBound {
		fmt.Printf("PVC %s is bound\n", pvc.Name)
		ready = true
	} else {
		fmt.Printf("PVC %s is not yet bound\n", pvc.Name)
	}
	return ready
}

func ToVolumeResource(pvc *apiv1.PersistentVolumeClaim) (*Resource, error) {
	return &Resource{
		name:  pvc.Name,
		ready: IsVolumeReady(pvc),
	}, nil
}

func (v *VolumeClientSet) CreateResource() error {
	result, err := v.clientSet.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault).Create(context.TODO(), v.pvc, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Volume %s created\n", result.Name)
	return nil
}

func (v *VolumeClientSet) GetResource() (*Resource, error) {
	result, err := v.clientSet.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault).Get(context.TODO(), v.pvc.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	v.pvc = result
	return ToVolumeResource(result)
}

func (v *VolumeClientSet) UpdateResource() error {
	fmt.Println("Volume update is not implemented")
	return nil
}

func (v *VolumeClientSet) DeleteResource() error {
	err := v.clientSet.CoreV1().PersistentVolumeClaims(apiv1.NamespaceDefault).Delete(context.TODO(), v.pvc.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Volume deleted")
	return nil
}

type DeploymentControllerClient struct {
	client     client.Client
	deployment *appsv1.Deployment
}

func (d *DeploymentControllerClient) CreateResource() error {
	return d.client.Create(context.TODO(), d.deployment)
}

func (d *DeploymentControllerClient) GetResource() (*Resource, error) {
	deployment := &appsv1.Deployment{}
	err := d.client.Get(context.TODO(), client.ObjectKey{
		Namespace: apiv1.NamespaceDefault,
		Name:      d.deployment.Name,
	}, deployment)
	if err != nil {
		return nil, err
	}
	d.deployment = deployment
	return ToDeploymentResource(deployment)
}

func (d *DeploymentControllerClient) UpdateResource() error {
	deployment := &appsv1.Deployment{}
	err := d.client.Get(context.TODO(), client.ObjectKey{
		Namespace: apiv1.NamespaceDefault,
		Name:      d.deployment.Name,
	}, deployment)
	if err != nil {
		return err
	}
	deployment.Spec.Replicas = int32Ptr(1)
	err = d.client.Update(context.TODO(), deployment)
	if err != nil {
		return err
	}
	return nil
}

func (d *DeploymentControllerClient) DeleteResource() error {
	return d.client.Delete(context.TODO(), d.deployment)
}

type ServiceControllerClient struct {
	client  client.Client
	service *apiv1.Service
}

func (s *ServiceControllerClient) CreateResource() error {
	return s.client.Create(context.TODO(), s.service)
}

func (s *ServiceControllerClient) GetResource() (*Resource, error) {
	service := &apiv1.Service{}
	err := s.client.Get(context.TODO(), client.ObjectKey{
		Namespace: apiv1.NamespaceDefault,
		Name:      s.service.Name,
	}, service)
	if err != nil {
		return nil, err
	}
	s.service = service
	return ToServiceResource(service)
}

func (s *ServiceControllerClient) UpdateResource() error {
	fmt.Println("Service update is not implemented")
	return nil
}

func (s *ServiceControllerClient) DeleteResource() error {
	return s.client.Delete(context.TODO(), s.service)
}

type VolumeControllerClient struct {
	client client.Client
	pvc    *apiv1.PersistentVolumeClaim
}

func (v *VolumeControllerClient) CreateResource() error {
	return v.client.Create(context.TODO(), v.pvc)
}

func (v *VolumeControllerClient) GetResource() (*Resource, error) {
	pvc := &apiv1.PersistentVolumeClaim{}
	err := v.client.Get(context.TODO(), client.ObjectKey{
		Namespace: apiv1.NamespaceDefault,
		Name:      v.pvc.Name,
	}, pvc)
	if err != nil {
		return nil, err
	}
	v.pvc = pvc
	return ToVolumeResource(pvc)
}

func (v *VolumeControllerClient) UpdateResource() error {
	fmt.Println("Volume update is not implemented")
	return nil
}

func (v *VolumeControllerClient) DeleteResource() error {
	return v.client.Delete(context.TODO(), v.pvc)
}

type ConfigType int

const (
	InCluster ConfigType = iota
	BuildConfig
)

func GetConfig(configType ConfigType) (*rest.Config, error) {
	switch configType {
	case BuildConfig:
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		return clientcmd.BuildConfigFromFlags("", *kubeconfig)
	case InCluster:
		return rest.InClusterConfig()
	default:
		return nil, errors.New("invalid config type")
	}
}

func main() {
	// Toggle between client-go and controller-runtime based on environment variable
	useControllerRuntime := os.Getenv("USE_CONTROLLER_RUNTIME") == "true"

	configType := BuildConfig // or InCluster, depending on environment
	config, err := GetConfig(configType)
	if err != nil {
		log.Fatalf("Failed to get kube config: %v", err)
	}

	var clients []K8sClient

	// Initialize objects for Deployment, Service, and Volume
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-deployment",
			Namespace: apiv1.NamespaceDefault,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
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

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-service",
			Namespace: apiv1.NamespaceDefault,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt32(80),
				},
			},
			Type: apiv1.ServiceTypeClusterIP,
		},
	}

	pvc := &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-pvc",
			Namespace: apiv1.NamespaceDefault,
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []apiv1.PersistentVolumeAccessMode{
				apiv1.ReadWriteOnce,
			},
			Resources: apiv1.VolumeResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	if useControllerRuntime {
		// Initialize controller-runtime based client
		fmt.Println("Using controller-runtime client")

		// Controller-runtime client initialization
		ctrlClient, err := client.New(config, client.Options{})
		if err != nil {
			log.Fatalf("Failed to create controller-runtime client: %v", err)
		}

		// Add controller-runtime clients to the slice
		clients = append(clients, &DeploymentControllerClient{client: ctrlClient, deployment: deployment})
		clients = append(clients, &ServiceControllerClient{client: ctrlClient, service: service})
		clients = append(clients, &VolumeControllerClient{client: ctrlClient, pvc: pvc})

	} else {
		// Initialize client-go based client
		fmt.Println("Using client-go client")
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Failed to create clientset: %v", err)
		}

		// Add client-go clients to the slice
		clients = append(clients, &DeploymentClientSet{clientSet: clientset, deployment: deployment})
		clients = append(clients, &ServiceClientSet{clientSet: clientset, service: service})
		clients = append(clients, &VolumeClientSet{clientSet: clientset, pvc: pvc})
	}

	// Perform CRUD operations in sequence
	for _, k8sClient := range clients {
		// Create Resource
		err := k8sClient.CreateResource()
		if err != nil {
			log.Fatalf("Failed to create resource: %v", err)
		}
	}

	for _, k8sClient := range clients {
		err = waitForResourceToBeReady(k8sClient, 10)
		if err != nil {
			log.Fatal(err)
		}

		// Update Resource
		err = k8sClient.UpdateResource()
		if err != nil {
			log.Fatalf("Failed to update resource: %v", err)
		}
	}

	for _, k8sClient := range clients {
		err = waitForResourceToBeReady(k8sClient, 10)
		if err != nil {
			log.Fatal(err)
		}
		err = k8sClient.DeleteResource()
		if err != nil {
			log.Fatalf("Failed to delete resource: %v", err)
		}
	}
}

func waitForResourceToBeReady(client K8sClient, retries int) error {
	for retries > 0 {
		res, err := client.GetResource()
		if err != nil {
			return err
		}

		if res.ready {
			return nil
		}

		fmt.Printf("Checking if resource %s is ready...", res.name)
		time.Sleep(10 * time.Second) // Adjust the interval as necessary
		retries--
	}
	return errors.New("retry Time Out, resource create/update took longer than expected time")
}
