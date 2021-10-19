
package main

import (
	"assignment1/configmap"
	"assignment1/deployments"
	"assignment1/services"
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crtconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	controllerClient, err := client.New(crtconfig.GetConfigOrDie(),client.Options{})
	if err != nil {
		fmt.Println("Failed to create rtc client", err)
	}
	rtc := client.NewNamespacedClient(controllerClient, apiv1.NamespaceDefault)

	fmt.Println("Running using the Client GO library to create DEPLOYMENT, SERVICES, CONFIGMAPS")
	deployments.CreateDeploymentClient(clientset)
	services.CreateServicesClient(clientset)
	configmap.CreateConfigMapClient(clientset)

	fmt.Println("Creating Resources")
	deployments.CreateDeployment()
	services.CreateServices()
	configmap.CreateConfigMap()

	time.Sleep(5 * time.Second)
	fmt.Println("Created Following Resources")
	deployments.GetAllDeployments()
	services.GetAllServices()
	configmap.GetAllConfigMaps()


	fmt.Println("Updating Resources")
	deployments.UpdateDeployment()
	services.UpdateServices()
	configmap.UpdateConfigMap()

	time.Sleep(5 * time.Second)
	fmt.Println("Deleting Resources")
	deployments.DeleteDeployment()
	services.DeleteService()
	configmap.DeleteConfigMap()

	time.Sleep(15 * time.Second)

	fmt.Println("Running using the controller runtime library to create DEPLOYMENT, SERVICES, CONFIGMAPS")
	deployments.SetRtClient(rtc)
	services.SetRtClient(rtc)
	configmap.SetRtClient(rtc)

	fmt.Println("Creating Resources")
	deployments.CreateRTCDeployment()
	services.CreateRTCServices()
	configmap.CreateRTCConfigMap()

	time.Sleep(5 * time.Second)
	fmt.Println("Created Following Resources")
	deployments.ListRTCDeployments()
	services.ListRTCServices()
	configmap.ListRTCConfigMaps()


	fmt.Println("Updating Resources")
	deployments.UpdateRTCDeployment()
	services.UpdateRTCService()
	configmap.UpdateRTCConfigMap()

	time.Sleep(5 * time.Second)
	fmt.Println("Deleting Resources")
	deployments.DeleteRTCDeployment()
	services.DeleteRTCService()
	configmap.DeleteRTCConfigMap()
}