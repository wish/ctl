package types

import (
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// CronJobDiscovery represents a cron job with the context information
type CronJobDiscovery struct {
	Context string
	batchv1beta1.CronJob
}

// GetLabels allows CronJobDiscovery to implement the Labeled interface
func (c CronJobDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// JobDiscovery represents a job with the context information
type JobDiscovery struct {
	Context string
	batchv1.Job
}

// GetLabels allows JobDiscovery to implement the Labeled interface
func (c JobDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// PodDiscovery represents a pod with the context information
type PodDiscovery struct {
	Context string
	corev1.Pod
}

// GetLabels allows PodDiscovery to implement the Labeled interface
func (c PodDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// ConfigMapDiscovery represents a configmap with the context information
type ConfigMapDiscovery struct {
	Context string
	corev1.ConfigMap
}

// GetLabels allows ConfigMapDiscovery to implement the Labeled interface
func (c ConfigMapDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// DeploymentDiscovery represents a deployment with the context information
type DeploymentDiscovery struct {
	Context string
	v1.Deployment
}

// GetLabels allows DeploymentDiscovery to implement the Labeled interface
func (c DeploymentDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// ReplicaSetDiscovery represents a deployment with the context information
type ReplicaSetDiscovery struct {
	Context string
	v1.ReplicaSet
}

// GetLabels allows ReplicaSetDiscovery to implement the Labeled interface
func (c ReplicaSetDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// RunDetails represents the commands and manifest to apply to launch an adhoc pod
type RunDetails struct {
	Resources    resource `json:"resources"`
	Active       bool     `json:"active"`
	Manifest     string   `json:"manifest"`
	PreLogin	 [][]string `json:"pre_login_command,omitempty"`
	LoginCommand []string `json:"login_command"`
}

type resource struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// ManifestDetails represents the manifest details of a running adhoc pod to get attributes like the namespace
type ManifestDetails struct {
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	APIVersion string `json:"apiVersion"`
	Kind 	   string `json:"kind"`
	Spec 	   struct {
		ActiveDeadlineSeconds string `json:"activeDeadlineSeconds"`
		Template	template `json:"template"`
	} `json:"spec"`

}

type template struct {
	Spec 		struct {
		Containers    []container `json:"containers"`
	}  `json:"spec"`
}

type container struct {
	Name	string `json:"name"`
	Command	[]string `json:"command,omitempty"`
	Args	[]string `json:"args,omitempty"`
	Image 	string `json:"image"`
}
