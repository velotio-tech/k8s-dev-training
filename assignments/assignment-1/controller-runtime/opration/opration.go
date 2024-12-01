package opration

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var clientset client.Client

var (
	deployemntName         = "nginx-deployment"
	labelkey               = "tier"
	labelvalue             = "frontend"
	image           string = "nginx"
	port            int32  = 80
	replicas        int32  = 2
	updatereplicase int32  = 5

	serviceName = "nginx-service"

	configmapName                   = "nginx-configmap"
	configData    map[string]string = map[string]string{
		"validatior": "manager",
	}
	configDataUpdated map[string]string = map[string]string{
		"validatior": "super-manager",
	}
)

func New() {
	env := os.Getenv("env")
	configpath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if env == "prod" {
		configpath = ""
	}
	config, err := clientcmd.BuildConfigFromFlags("", configpath)
	if err != nil {
		panic(err)
	}

	k8client, err := client.New(config, client.Options{})
	if err != nil {
		panic(err)
	}
	clientset = client.NewNamespacedClient(k8client, "default")
}

func Deployment() {
	for {
		var opration int
		fmt.Printf("Choose deployment opration\n 1) Create\n 2) Read\n 3) Update\n 4) Delete\n 5) Exit \n Input : ")
		fmt.Scanf("%d", &opration)

		switch opration {
		case 1:
			createDeployment()
		case 2:
			readDeployment()
		case 3:
			updateDeployment()
		case 4:
			deleteDeployment()
		case 5:
			break
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func Service() {
	for {
		var opration int
		fmt.Printf("Choose service opration\n 1) Create\n 2) Read\n 3) Update\n 4) Delete\n 5) Exit \n Input : ")
		fmt.Scanf("%d", &opration)

		switch opration {
		case 1:
			createService()
		case 2:
			readService()
		case 3:
			updateService()
		case 4:
			deleteService()
		case 5:
			break
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func ConfigMap() {
	for {
		var opration int
		fmt.Printf("Choose configmap opration\n 1) Create\n 2) Read\n 3) Update\n 4) Delete\n 5) Exit \n Input : ")
		fmt.Scanf("%d", &opration)

		switch opration {
		case 1:
			createConfigMap()
		case 2:
			readConfigMap()
		case 3:
			updateConfigMap()
		case 4:
			deleteConfigMap()
		case 5:
			break
		default:
			fmt.Println("Invalid choice")
		}
	}
}
