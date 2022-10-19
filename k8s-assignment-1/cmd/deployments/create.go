/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package deployments

import (
	"context"
	"fmt"

	cmdv1 "github.com/apapapap/k8s-dev-training/assignment-1/kube-client/cmd"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createCmd represents the create deployment command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create deployment",
	Long:  "kube-client create deployment",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		namespace := "default"
		deployment := getDeploymentObj()
		if cmdv1.UseCtrlRuntime {
			err = cmdv1.CtrlClient.Create(context.Background(), deployment)
		} else {
			_, err = cmdv1.ClientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
		}
		if err != nil {
			fmt.Println("Error creating deployment, error: ", err)
			return
		}
		fmt.Println("Deployment created successfully")
	},
}

func init() {
	deploymentsCmd.AddCommand(createCmd)
}

func getDeploymentObj() *appsv1.Deployment {
	return &appsv1.Deployment{
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
}

func int32Ptr(i int32) *int32 { return &i }
