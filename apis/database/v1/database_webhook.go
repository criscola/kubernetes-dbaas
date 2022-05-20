/*
Copyright 2021.

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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var databaselog = logf.Log.WithName("database-resource-webhook")

func (r *Database) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-database-dbaas-bedag-ch-v1-database,mutating=true,failurePolicy=fail,sideEffects=None,groups=database.dbaas.bedag.ch,resources=databases,verbs=create;update,versions=v1,name=mdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Database{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Database) Default() {
	//databaselog.Info("default", "name", r.Name)

	// TODO(user): fill in with defaulting logic
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-database-dbaas-bedag-ch-v1-database,mutating=false,failurePolicy=fail,sideEffects=None,groups=database.dbaas.bedag.ch,resources=databases,verbs=create;update,versions=v1,name=vdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Database{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateCreate() error {
	return nil
}

// ValidateUpdate disables any update to the 'spec' field of Database resources.
func (r *Database) ValidateUpdate(old runtime.Object) error {
	databaselog.Info("validate update", "name", r.Name)
	var allErrs field.ErrorList

	rOld := old.(*Database)

	databaselog.Info("validate update", "oldSpec", rOld.Spec)
	databaselog.Info("validate update", "newSpec", r.Spec)

	if !reflect.DeepEqual(r.Spec, rOld.Spec) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), r.Spec, "update operations not allowed, please explicitly "+
			"delete the resource in order to recreate it."))

		return apierrors.NewInvalid(schema.GroupKind{Group: "database.dbaas.bedag.ch", Kind: "Database"},
			"database", allErrs)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateDelete() error {
	return nil
}
