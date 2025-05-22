/*


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

package v1beta1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var etcdadmconfiglog = logf.Log.WithName("etcdadmconfig-resource")

func (r *EtcdadmConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithDefaulter(r).
		WithValidator(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-bootstrap-cluster-x-k8s-io-v1beta1-etcdadmconfig,mutating=true,failurePolicy=fail,groups=bootstrap.cluster.x-k8s.io,resources=etcdadmconfigs,verbs=create;update,versions=v1beta1,name=metcdadmconfig.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &EtcdadmConfig{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type
func (r *EtcdadmConfig) Default(_ context.Context, obj runtime.Object) error {
	etcdadmConfig, ok := obj.(*EtcdadmConfig)
	if !ok {
		return fmt.Errorf("expected an EtcdadmConfig but got %T", obj)
	}

	etcdadmconfiglog.Info("default", "name", etcdadmConfig.Name)

	// TODO(user): fill in your defaulting logic.
	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-bootstrap-cluster-x-k8s-io-v1beta1-etcdadmconfig,mutating=false,failurePolicy=fail,groups=bootstrap.cluster.x-k8s.io,resources=etcdadmconfigs,versions=v1beta1,name=vetcdadmconfig.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomValidator = &EtcdadmConfig{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type
func (r *EtcdadmConfig) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	etcdadmConfig, ok := obj.(*EtcdadmConfig)
	if !ok {
		return nil, fmt.Errorf("expected an EtcdadmConfig but got %T", obj)
	}

	etcdadmconfiglog.Info("validate create", "name", etcdadmConfig.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type
func (r *EtcdadmConfig) ValidateUpdate(_ context.Context, old, obj runtime.Object) (admission.Warnings, error) {
	etcdadmConfig, ok := obj.(*EtcdadmConfig)
	if !ok {
		return nil, fmt.Errorf("expected an EtcdadmConfig but got %T", obj)
	}

	etcdadmconfiglog.Info("validate update", "name", etcdadmConfig.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type
func (r *EtcdadmConfig) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	etcdadmConfig, ok := obj.(*EtcdadmConfig)
	if !ok {
		return nil, fmt.Errorf("expected an EtcdadmConfig but got %T", obj)
	}

	etcdadmconfiglog.Info("validate delete", "name", etcdadmConfig.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
