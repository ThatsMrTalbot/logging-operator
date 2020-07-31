// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var flowlog = logf.Log.WithName("flow-resource")

func (r *Flow) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-logging-banzaicloud-io-banzaicloud-io-v1beta1-flow,mutating=false,failurePolicy=fail,groups=logging.banzaicloud.io.banzaicloud.io,resources=flows,versions=v1beta1,name=vflow.kb.io

var _ webhook.Validator = &Flow{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Flow) ValidateCreate() error {
	flowlog.Info("validate create", "name", r.Name)
	return r.ValidateFields().ToAggregate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Flow) ValidateUpdate(old runtime.Object) error {
	flowlog.Info("validate update", "name", r.Name)
	return r.ValidateFields().ToAggregate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Flow) ValidateDelete() error {
	// No deletion validation implemented
	return nil
}

// ValidateFields performs validation on each field of the object, it does not go as far
// as to load the contents of secrets and validate the contents as this is to be called
// as part of a webhook.
func (r *Flow) ValidateFields() field.ErrorList {
	// We only validate the spec
	return validateFlowSpec(field.NewPath("spec"), r.Spec)
}

func validateFlowSpec(path *field.Path, spec FlowSpec) field.ErrorList {
	// Create variable to track validation errors
	allErrs := field.ErrorList{}

	// Validate filters
	allErrs = append(allErrs, validateFilters(path.Child("filters"), spec.Filters)...)

	// Return any errors
	return allErrs
}

func validateFilters(path *field.Path, filters []Filter) field.ErrorList {
	// Create variable to track validation errors
	allErrs := field.ErrorList{}

	// Iterate over filters, validating them
	for i, filter := range filters {
		allErrs = append(allErrs, validateFilter(path.Index(i), filter)...)
	}

	// Return any errors
	return allErrs
}

func validateFilter(path *field.Path, filter Filter) field.ErrorList {
	// Create variable to track validation errors
	allErrs := field.ErrorList{}

	// Validate filters
	allErrs = append(allErrs, filter.StdOut.ValidateFields(path.Child("stdout"))...)
	allErrs = append(allErrs, filter.Parser.ValidateFields(path.Child("parser"))...)
	// allErrs = append(allErrs, filter.TagNormaliser.ValidateFields(path.Child("tag_normaliser"))...)
	// allErrs = append(allErrs, filter.DedotFilterConfig.ValidateFields(path.Child("dedot"))...)
	// allErrs = append(allErrs, filter.RecordTransformer.ValidateFields(path.Child("record_transformer"))...)
	// allErrs = append(allErrs, filter.RecordModifier.ValidateFields(path.Child("record_modifier"))...)
	// allErrs = append(allErrs, filter.GeoIP.ValidateFields(path.Child("geoip"))...)
	// allErrs = append(allErrs, filter.Concat.ValidateFields(path.Child("concat"))...)
	// allErrs = append(allErrs, filter.DetectExceptions.ValidateFields(path.Child("detectExceptions"))...)
	// allErrs = append(allErrs, filter.Grep.ValidateFields(path.Child("grep"))...)
	// allErrs = append(allErrs, filter.PrometheusConfig.ValidateFields(path.Child("prometheus"))...)
	// allErrs = append(allErrs, filter.Throttle.ValidateFields(path.Child("throttle"))...)

	// Return any errors
	return allErrs
}
