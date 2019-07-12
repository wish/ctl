package types

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// Found resources with the context as client-go does not contain context
type CronJobDiscovery struct {
	Context string
	batchv1beta1.CronJob
}

func (c CronJobDiscovery) GetLabels() map[string]string {
	return c.Labels
}

type RunDiscovery struct {
	Context string
	batchv1.Job
}

func (c RunDiscovery) GetLabels() map[string]string {
	return c.Labels
}

type PodDiscovery struct {
	Context string
	corev1.Pod
}

func (c PodDiscovery) GetLabels() map[string]string {
	return c.Labels
}
