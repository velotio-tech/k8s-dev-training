package runtimectrl

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pankaj9310/k8s-dev-training/pankaj/assignment1/goclient"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var roleName = "demo-role"

//RoleOperations - perform CURD operations on role object using controller runtime library
func RoleOperations() {

	ctrClient := getClient()

	fileBytes, err := ioutil.ReadFile("configfile/role.yaml")
	if err != nil {
		panic(err.Error())
	}

	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(string(fileBytes))), 1024)

	var roleSpec rbacv1.Role
	err = dec.Decode(&roleSpec)

	if err != nil {
		panic(err.Error())
	}

	// Create Role
	fmt.Println("Role create operation start ........")
	err = ctrClient.Create(context.Background(), &roleSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Role Create operation completed ........")

	//Read Role
	fmt.Println("Role Read operation start ........")
	err = ctrClient.Get(context.Background(), client.ObjectKey{Namespace: namespace, Name: roleName}, &roleSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Role name: %q, in namespace %s, Role.Rules[0].Resources accessible %q.\n", roleSpec.Name, namespace, roleSpec.Rules[0].Resources)
	fmt.Println("Role Read operation completed ........")

	//Update Role
	updateResources := []string{"pods", "services"}
	fmt.Println("Role Update operation start ........")
	roleSpec.Rules[0].Resources = updateResources
	err = ctrClient.Update(context.Background(), &roleSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated: Role name: %q, in namespace %s, Role.Rules[0].Resources accessible %q.\n", roleSpec.Name, namespace, roleSpec.Rules[0].Resources)
	fmt.Println("Role Update operation completed ........")

	//Verify Role
	fmt.Println("Role Update verfication operation start ........")

	err = ctrClient.Update(context.Background(), &roleSpec)
	if err != nil {

		panic(err.Error())
	}
	if goclient.ComapreSlice(updateResources, roleSpec.Rules[0].Resources) {
		fmt.Println("Role Verified Successfully")
	} else {
		fmt.Printf("Role Verfication failed, expected role %q; Actual Role is  %q.\n", updateResources, roleSpec.Rules[0].Resources)
	}

	fmt.Println("Role Update verfication operation completed ........")

	// Delete Deployments
	fmt.Println("Role Delete operation started ........")
	err = ctrClient.Delete(context.Background(), &roleSpec)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Role Delete operation completed ........")
}
