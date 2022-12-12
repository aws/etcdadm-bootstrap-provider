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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capbk "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
const (
	DataSecretAvailableCondition clusterv1.ConditionType = "DataSecretAvailable"
	// CloudConfig make the bootstrap data to be of cloud-config format.
	CloudConfig Format = "cloud-config"
	// Bottlerocket make the bootstrap data to be of bottlerocket format.
	Bottlerocket Format = "bottlerocket"
)

// Format specifies the output format of the bootstrap data
// +kubebuilder:validation:Enum=cloud-config;bottlerocket
type Format string

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EtcdadmConfigSpec defines the desired state of EtcdadmConfig
type EtcdadmConfigSpec struct {
	// Users specifies extra users to add
	// +optional
	Users []capbk.User `json:"users,omitempty"`

	// +optional
	EtcdadmBuiltin bool `json:"etcdadmBuiltin,omitempty"`

	// +optional
	EtcdadmInstallCommands []string `json:"etcdadmInstallCommands,omitempty"`

	// PreEtcdadmCommands specifies extra commands to run before kubeadm runs
	// +optional
	PreEtcdadmCommands []string `json:"preEtcdadmCommands,omitempty"`

	// PostEtcdadmCommands specifies extra commands to run after kubeadm runs
	// +optional
	PostEtcdadmCommands []string `json:"postEtcdadmCommands,omitempty"`

	// Format specifies the output format of the bootstrap data
	// +optional
	Format Format `json:"format,omitempty"`

	// BottlerocketConfig specifies the configuration for the bottlerocket bootstrap data
	// +optional
	BottlerocketConfig *BottlerocketConfig `json:"bottlerocketConfig,omitempty"`

	// CloudInitConfig specifies the configuration for the cloud-init bootstrap data
	// +optional
	CloudInitConfig *CloudInitConfig `json:"cloudInitConfig,omitempty"`

	// Files specifies extra files to be passed to user_data upon creation.
	// +optional
	Files []capbk.File `json:"files,omitempty"`

	// Proxy holds the https and no proxy information
	// This is only used for bottlerocket
	// +optional
	Proxy *ProxyConfiguration `json:"proxy,omitempty"`

	// RegistryMirror holds the image registry mirror information
	// This is only used for bottlerocket
	// +optional
	RegistryMirror *RegistryMirrorConfiguration `json:"registryMirror,omitempty"`

	// CipherSuites is a list of comma-delimited supported TLS cipher suites, mapping to the --cipher-suites flag.
	// Default is empty, which means that they will be auto-populated by Go.
	// +optional
	CipherSuites string `json:"cipherSuites,omitempty"`
}

type BottlerocketConfig struct {
	// EtcdImage specifies the etcd image to use by etcdadm
	EtcdImage string `json:"etcdImage,omitempty"`

	// BootstrapImage specifies the container image to use for bottlerocket's bootstrapping
	BootstrapImage string `json:"bootstrapImage"`

	// AdminImage specifies the admin container image to use for bottlerocket.
	AdminImage string `json:"adminImage"`

	// ControlImage specifies the control container image to use for bottlerocket.
	ControlImage string `json:"controlImage"`

	// PauseImage specifies the image to use for the pause container
	PauseImage string `json:"pauseImage"`

	// CustomHostContainers adds additional host containers for bottlerocket.
	// +optional
	CustomHostContainers []BottlerocketHostContainer `json:"customHostContainers,omitempty"`

	// CustomBootstrapContainers adds additional bootstrap containers for bottlerocket.
	// +optional
	CustomBootstrapContainers []BottlerocketBootstrapContainer `json:"customBootstrapContainers,omitempty"`
}

// BottlerocketHostContainer holds the host container setting for bottlerocket.
type BottlerocketHostContainer struct {
	// Name is the host container name that will be given to the container in BR's `apiserver`
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Superpowered indicates if the container will be superpowered
	// +kubebuilder:validation:Required
	Superpowered bool `json:"superpowered"`

	// Image is the actual location of the host container image.
	Image string `json:"image"`

	// UserData is the userdata that will be attached to the image.
	// +optional
	UserData string `json:"userData,omitempty"`
}

// BottlerocketBootstrapContainer holds the bootstrap container setting for bottlerocket.
type BottlerocketBootstrapContainer struct {
	// Name is the bootstrap container name that will be given to the container in BR's `apiserver`.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Image is the actual image used for Bottlerocket bootstrap.
	Image string `json:"image"`

	// Essential decides whether or not the container should fail the boot process.
	// Bootstrap containers configured with essential = true will stop the boot process if they exit code is a non-zero value.
	// Default is false.
	// +optional
	Essential bool `json:"essential"`

	// Mode represents the bootstrap container mode.
	// +kubebuilder:validation:Enum=always;off;once
	Mode string `json:"mode"`

	// UserData is the base64-encoded userdata.
	// +optional
	UserData string `json:"userData,omitempty"`
}

type CloudInitConfig struct {
	// +optional
	Version string `json:"version,omitempty"`

	// EtcdReleaseURL is an optional field to specify where etcdadm can download etcd from
	// +optional
	EtcdReleaseURL string `json:"etcdReleaseURL,omitempty"`

	// InstallDir is an optional field to specify where etcdadm will extract etcd binaries to
	// +optional
	InstallDir string `json:"installDir,omitempty"`
}

// ProxyConfiguration holds the settings for proxying bottlerocket services
type ProxyConfiguration struct {
	// HTTP Proxy
	HTTPProxy string `json:"httpProxy,omitempty"`

	// HTTPS proxy
	HTTPSProxy string `json:"httpsProxy,omitempty"`

	// No proxy, list of ips that should not use proxy
	NoProxy []string `json:"noProxy,omitempty"`
}

// RegistryMirrorConfiguration holds the settings for image registry mirror
type RegistryMirrorConfiguration struct {
	// Endpoint defines the registry mirror endpoint to use for pulling images
	Endpoint string `json:"endpoint,omitempty"`

	// CACert defines the CA cert for the registry mirror
	CACert string `json:"caCert,omitempty"`
}

// EtcdadmConfigStatus defines the observed state of EtcdadmConfig
type EtcdadmConfigStatus struct {
	// Conditions defines current service state of the KubeadmConfig.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	DataSecretName *string `json:"dataSecretName,omitempty"`

	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// EtcdadmConfig is the Schema for the etcdadmconfigs API
type EtcdadmConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdadmConfigSpec   `json:"spec,omitempty"`
	Status EtcdadmConfigStatus `json:"status,omitempty"`
}

func (e *EtcdadmConfig) GetConditions() clusterv1.Conditions {
	return e.Status.Conditions
}

func (e *EtcdadmConfig) SetConditions(conditions clusterv1.Conditions) {
	e.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// EtcdadmConfigList contains a list of EtcdadmConfig
type EtcdadmConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdadmConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdadmConfig{}, &EtcdadmConfigList{})
}
