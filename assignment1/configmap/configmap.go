package configmap

import (
  "context"
  "fmt"
  apiv1 "k8s.io/api/core/v1"
  v1 "k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
  "k8s.io/client-go/util/retry"
  "log"
)

var configMapClient corev1.ConfigMapInterface

func CreateConfigMapClient(clientset *kubernetes.Clientset)  {
  configMapClient = clientset.CoreV1().ConfigMaps(apiv1.NamespaceDefault)
}

func GetAllConfigMaps() {
  fmt.Printf("Listing ConfigMaps in namespace %q:\n", apiv1.NamespaceDefault)
  list, err := configMapClient.List(context.TODO(), metav1.ListOptions{})
  if err != nil {
    panic(err)
  }
  for _, d := range list.Items {
    fmt.Printf(" * %s \n", d.Name)
  }
}

func CreateConfigMap() {
  configMapData := make(map[string]string, 0)
  uiProperties := `
  color.good=purple
  color.bad=yellow
  allow.textmode=true`
  configMapData["ui.properties"] = uiProperties

  var configMap = &v1.ConfigMap{
    TypeMeta: metav1.TypeMeta{
      Kind:       "ConfigMap",
      APIVersion: "v1",
    },
    ObjectMeta: metav1.ObjectMeta{
      Name:      "my-data",
    },
    Data: configMapData,
  }
  fmt.Println("Creating CM...")
  result, err := configMapClient.Create(context.TODO(), configMap, metav1.CreateOptions{})
  if err != nil {
    log.Println("Error Occurred while creating the configMap")
  }
  fmt.Printf("Created ConfigMap %q.\n", result.GetObjectMeta().GetName())
}

func DeleteConfigMap(){
  fmt.Println("Deleting configMap...")
  deletePolicy := metav1.DeletePropagationForeground
  if err := configMapClient.Delete(context.TODO(), "my-data", metav1.DeleteOptions{
    PropagationPolicy: &deletePolicy,
  }); err != nil {
    panic(err)
  }
  fmt.Println("Deleted CM.")
}

func UpdateConfigMap() {
  retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
    result, getErr := configMapClient.Get(context.TODO(), "my-data", metav1.GetOptions{})
    if getErr != nil {
      log.Println(fmt.Errorf("Failed to get latest version of ConfigMap: %v", getErr))
    }

    result.Data["ui.properties.color.bad"] = "red"

    _, updateErr := configMapClient.Update(context.TODO(), result, metav1.UpdateOptions{})
    return updateErr
  })
  if retryErr != nil {
    panic(fmt.Errorf("Update failed: %v", retryErr))
  }
  fmt.Println("Updated CM...")
}
