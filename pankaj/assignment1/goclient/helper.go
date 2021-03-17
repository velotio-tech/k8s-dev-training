package goclient

import (
	"path/filepath"
	"sort"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func KubeConfig() *kubernetes.Clientset {

	var kubernetesConfig *rest.Config
	homeDir := homedir.HomeDir()
	var err error

	if homeDir != "" {
		kubernetesConfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir, ".kube", "config"))
		if err != nil {
			panic(err)
		}
	} else {
		kubernetesConfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func ConvertToPtr(val int32) *int32 {
	return &val

}

func ComapreSlice(str1 []string, str2 []string) bool {
	sort.Strings(str1)
	sort.Strings(str2)
	for idx, val := range str1 {
		if str2[idx] != val {
			return false
		}
	}
	return true
}
