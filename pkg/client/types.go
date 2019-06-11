package client

import (
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// Found resources with the context as client-go does not contain context
type CronJobDiscovery struct {
	Context string
	batchv1beta1.CronJob
}

type PodDiscovery struct {
	Context string
	corev1.Pod
}
