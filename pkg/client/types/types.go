package types

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
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
	extensionsv1beta1.Deployment
}

// GetLabels allows DeploymentDiscovery to implement the Labeled interface
func (c DeploymentDiscovery) GetLabels() map[string]string {
	return c.Labels
}

// ReplicaSetDiscovery represents a deployment with the context information
type ReplicaSetDiscovery struct {
	Context string
	extensionsv1beta1.ReplicaSet
}

// GetLabels allows ReplicaSetDiscovery to implement the Labeled interface
func (c ReplicaSetDiscovery) GetLabels() map[string]string {
	return c.Labels
}
