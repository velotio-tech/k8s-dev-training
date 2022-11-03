/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var mongodblog = logf.Log.WithName("mongodb-resource")

func (r *MongoDB) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-db-velotio-com-v1-mongodb,mutating=true,failurePolicy=fail,sideEffects=None,groups=db.velotio.com,resources=mongodbs,verbs=create;update,versions=v1,name=mmongodb.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &MongoDB{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MongoDB) Default() {
	mongodblog.Info("default", "name", r.Name)
	if r.Spec.MaxUsers < 1 {
		r.Spec.MaxUsers = 1
	}
	if r.Spec.MaxConcurrentConnections < 1 {
		r.Spec.MaxConcurrentConnections = 100
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-db-velotio-com-v1-mongodb,mutating=false,failurePolicy=fail,sideEffects=None,groups=db.velotio.com,resources=mongodbs,verbs=create;update,versions=v1,name=vmongodb.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MongoDB{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MongoDB) ValidateCreate() error {
	mongodblog.Info("validate create", "name", r.Name)
	if strings.TrimSpace(r.Spec.InitPassword) == "" || strings.TrimSpace(r.Spec.InitUser) == "" {
		mongodblog.Info("init user OR init password validation failed")
		return errors.New("init user or init password cannot be empty")
	}
	if r.Spec.GVK.Kind != "Pod" {
		mongodblog.Info(fmt.Sprintf("kind %s not supported", r.Spec.GVK.Kind))
		return errors.New("unsupported kind")
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MongoDB) ValidateUpdate(old runtime.Object) error {
	mongodblog.Info("validate update", "name", r.Name)
	oldObject := old.(*MongoDB)
	if oldObject.Spec.InitUser != r.Spec.InitUser || oldObject.Spec.InitPassword != r.Spec.InitPassword {
		mongodblog.Info("attemplting to update init user OR init password, validation failed")
		return fmt.Errorf("init user or init password cannot be updated")
	}
	if r.Spec.GVK.Kind != "Pod" || r.Spec.GVK.APIVersion != "v1" {
		mongodblog.Info(fmt.Sprintf("kind %s OR APIVerions %s not supported", r.Spec.GVK.Kind, r.Spec.GVK.APIVersion))
		return errors.New("unsupported kind")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MongoDB) ValidateDelete() error {
	mongodblog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
