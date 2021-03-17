package goclient

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var roleName = "demo-role"

//RoleOperations - perform CURD operations on Role object using client-go library
func RoleOperations() {
	clientset := KubeConfig()
	roleClient := clientset.RbacV1().Roles(namespace)

	fileBytes, err := ioutil.ReadFile("configfile/role.yaml")
	if err != nil {
		panic(err.Error())
	}

	var roleSpec rbacv1.Role
	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1024)

	err = dec.Decode(&roleSpec)

	if err != nil {
		panic(err.Error())
	}

	// Create Role
	fmt.Println("Role create operation start ........")
	role, err := roleClient.Create(context.TODO(), &roleSpec, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Create Role %q in namespace %s\n", role.GetObjectMeta().GetName(), namespace)
	fmt.Println("Role Create operation completed ........")

	//Get Role
	fmt.Println("Role Get operation start ........")
	role, err = roleClient.Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Role name: %q, in namespace %s, Role.Rules[0].Resources accessible %q.\n", role.GetObjectMeta().GetName(), namespace, role.Rules[0].Resources)
	fmt.Println("Role Get operation completed ........")

	//Update Role
	fmt.Println("Role Update operation start ........")
	updateResources := []string{"pods", "services"}
	role.Rules[0].Resources = updateResources
	role, err = roleClient.Update(context.TODO(), role, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated: Role name: %q, in namespace %s, Role.Rules[0].Resources accessible %q.\n", role.GetObjectMeta().GetName(), namespace, role.Rules[0].Resources)
	fmt.Println("Role Update operation completed ........")

	//Verify Role
	fmt.Println("Role Update verfication operation start ........")
	role, err = roleClient.Get(context.TODO(), roleName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	if ComapreSlice(updateResources, role.Rules[0].Resources) {
		fmt.Println("Role Verified Successfully")
	} else {
		fmt.Printf("Role Verfication failed, expected role %q; Actual Role is  %q.\n", updateResources, role.Rules[0].Resources)
	}

	fmt.Println("Role Update verfication operation completed ........")

	// Delete Deployments
	fmt.Println("Role Delete operation started ........")
	deletePolicy := metav1.DeletePropagationForeground
	err = roleClient.Delete(context.TODO(), roleName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Role Delete operation completed ........")
}
