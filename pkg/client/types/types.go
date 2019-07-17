package types

import (
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

// RunDiscovery represents a job with the context information
type RunDiscovery struct {
	Context string
	batchv1.Job
}

// GetLabels allows RunDiscovery to implement the Labeled interface
func (c RunDiscovery) GetLabels() map[string]string {
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
