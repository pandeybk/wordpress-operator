/*
Copyright 2018 Pressinfra SRL

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindWordpress = "Wordpress"
)

// SecretRef represents a reference to a Secret
type SecretRef string

// URL represents a valid URL string
type URL string

// Domain represents a valid domain name
type Domain string

const (
	// CPU in cores for nginx (eg. 500m = .5 cores)
	ResourceNginxCPU corev1.ResourceName = "nginx/cpu"
	// Memory, in bytes for nginx. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceNginxMemory corev1.ResourceName = "nginx/memory"
	// CPU in cores for PHP (eg. 500m = .5 cores)
	ResourcePHPCPU corev1.ResourceName = "php/cpu"
	// Memory, in bytes for PHP. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourcePHPMemory corev1.ResourceName = "php/memory"
	// Number of PHP workers
	ResourcePHPWorkers corev1.ResourceName = "php/workers"
	// Memory, in bytes for PHP worker. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourcePHPWorkerMemory corev1.ResourceName = "php/worker-memory"
	// Maximum execution time of a PHP worker in seconds
	ResourcePHPMaxExecutionSeconds corev1.ResourceName = "php/max-execution-seconds"
	// Maximum request body size in bytes (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceIngressBodySize corev1.ResourceName = "ingress/max-body-size"
	// CPU in cores for running wp cli commands (eg. 500m = .5 cores)
	ResourceCLICPU corev1.ResourceName = "cli/cpu"
	// Memory, in bytes for running cli commands. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceCLIMemory corev1.ResourceName = "cli/memory"
)

// +k8s:openapi-gen=true

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Wordpress struct {
	// +k8s:openapi-gen=false
	metav1.TypeMeta `json:",inline"`
	// +k8s:openapi-gen=false
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec WordpressSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type WordpressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Wordpress `json:"items"`
}

type WordpressSpec struct {
	// Domains for this this site answers. The first item is set as the "main
	// domain" (WP_HOME and WP_SITEURL constants).
	Domains []Domain `json:"domains"`
	// TLSSecretRef a secret containing the TLS certificates for this site.
	// +optional
	TLSSecretRef SecretRef `json:"tlsSecretRef,omitempty"`
	// ContentVolumeSpec defines how wp-content is mounted.
	// +optional
	ContentVolumeSpec WordpressVolumeSpec `json:"contentVolumeSpec,omitempty"`
	// MediaVolumeSpec if specified, defines the volume for storing media
	// files.
	// +optional
	MediaVolumeSpec *WordpressVolumeSpec `json:"mediaVolumeSpec,omitempty"`
	// The secret name which contain credentials and customizations for this
	// WordPress site. The secret is mounted as a volume, and the following keys
	// get special treatment:
	// - wp-config.php
	//   Custom wp-config
	// - php.ini
	//   Contains custom php.ini definitions
	// - id_rsa
	//   Is the SSH key used to access the code repository
	// - netrc
	//   Is the .netrc file used for cloning the code repository. It can be used
	//   for granting access to repos over HTTP
	// - google_service_account.json
	//   Google Service Account key file, for accessing Google Cloud Services
	//   from within the WordPress site
	// - aws_credentials
	// - aws_config
	//   The ~/.aws/credentials and ~/.aws/config files, used for accessing AWS
	//   Services from within the WordPress site
	// - nginx-server.conf
	//   nginx customizations to include in nginx http {  } config block
	// - nginx-vhost.conf
	//   nginx customizations to include in nginx server {  } config block
	SecretRef SecretRef `json:"secretRef"`
	// List of environment variables to set in the PHP container.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []corev1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
	// Image is the docker image to use as basis for the execution environment
	// of this WordPress site.
	// +optional
	Image string `json:"image,omitempty"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// CLI driver to use for running wp cron, database upgrades:
	// * standalone (default)
	//   spawns a standalone pod when running a cli command
	// * inline
	//   uses kubectl to exec into a running pod, and executes the cli commands
	//   there. These jobs tend to fail more often and are more fragile, but
	//   they work in cases where standalone driver cannot be used (eg. content
	//   must be shred by pods running wordpress and pods running jobs)
	CLIDriver string `json:"cliDriver,omitempty"`
	// Compute Resources required by this Wordpress instance.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// If specified, the pod's scheduling constraints.
	// +optional
	Affinity corev1.Affinity `json:"affinity,omitempty"`
	// ServiceAccount to use for running pods.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// If specified apply these annotations to the Ingress resource created for
	// this Wordpress Site.
	// +optional
	IngressAnnotations map[string]string `json:"ingressAnnotations,omitempty"`
	// ServiceSpec is the specification for the service created for this WordPress Site.
	// +optional
	ServiceSpec *corev1.ServiceSpec `json:"serviceSpec,omitempty"`
}

type WordpressVolumeSpec struct {
	// ReadOnly specifies if the volume should be mounted read only.
	// WARNING: It still can be mounted read-write for initialization for
	// example.
	// Defaults to false.
	// +optional
	ReadOnly *bool `json:"readOnly,omitempty"`
	// EmptyDir to use if no PersistentVolumeClaim or HostPath is specified
	// +optional
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`
	// HostPath to use instead of a PersistentVolumeClaim.
	// +optional
	HostPath *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`
	// PersistentVolumeClaim to use. It has the highest level of precedence,
	// followed by HostPath and EmptyDir
	// +optional
	PersistentVolumeClaim *corev1.PersistentVolumeClaimSpec `json:"persistentVolumeClaim,omitempty"`
}