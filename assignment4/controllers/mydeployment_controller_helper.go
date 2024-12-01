package controllers

import (
	"context"

	"github.com/go-logr/logr"
	velotiov1 "github.com/pankaj9310/k8s-dev-training/pankaj/assignment4/api/v1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func buildDeployment(myDeployment velotiov1.MyDeployment) *apps.Deployment {
	deployment := apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            myDeployment.Spec.DeploymentName,
			Namespace:       myDeployment.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&myDeployment, velotiov1.GroupVersion.WithKind("MyDeployment"))},
		},
		Spec: apps.DeploymentSpec{
			Replicas: myDeployment.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"example-controller.jetstack.io/deployment-name": myDeployment.Spec.DeploymentName,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"example-controller.jetstack.io/deployment-name": myDeployment.Spec.DeploymentName,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  "pankaj-myjob",
							Image: "docker.io/pankaj9310/myjob:latest",
							Command: []string{
								"./app",
							},
						},
					},
				},
			},
		},
	}
	return &deployment
}

// cleanupOwnedResources will Delete any existing Deployment resources that
// were created for the given MyDeployment that no longer match the
// myDeployment.spec.deploymentName field.
func (r *MyDeploymentReconciler) cleanupOwnedResources(ctx context.Context, log logr.Logger, myDeployment *velotiov1.MyDeployment) error {
	log.Info("finding existing Deployments for MyDeployment resource")

	// List all deployment resources owned by this MyDeployment
	var deployments apps.DeploymentList
	if err := r.List(ctx, &deployments, client.InNamespace(myDeployment.Namespace), client.MatchingField(deploymentOwnerKey, myDeployment.Name)); err != nil {
		return err
	}

	deleted := 0
	for _, depl := range deployments.Items {
		if depl.Name == myDeployment.Spec.DeploymentName {
			// If this deployment's name matches the one on the MyDeployment resource
			// then do not delete it.
			continue
		}

		if err := r.Client.Delete(ctx, &depl); err != nil {
			log.Error(err, "failed to delete Deployment resource")
			return err
		}

		r.Recorder.Eventf(myDeployment, core.EventTypeNormal, "Deleted", "Deleted deployment %q", depl.Name)
		deleted++
	}

	log.Info("finished cleaning up old Deployment resources", "number_deleted", deleted)

	return nil
}
