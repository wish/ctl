package client

import (
	"github.com/wish/ctl/pkg/client/types"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	// describeversioned "k8s.io/kubectl/pkg/describe/versioned"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/describe"
	"strings"
)

// Helper to print contextual info
func (c *Client) describeContextInfo(context string) string {
	if c.extension.ClusterExt == nil {
		return ""
	}
	if _, ok := c.extension.ClusterExt[context]; !ok {
		return ""
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Context: \t%s\n", context)
	if _, ok := c.extension.ClusterExt[context]; ok {
		fmt.Fprintf(&sb, "Labels:  \t")
		first := true
		for k, v := range c.extension.ClusterExt[context] {
			if strings.HasPrefix(k, "_") {
				continue
			}
			if first {
				fmt.Fprintf(&sb, "%s=%s\n", k, v)
				first = false
			} else {
				fmt.Fprintf(&sb, "         \t%s=%s\n", k, v)
			}
		}
	}
	return sb.String()
}

// DescribePod returns a human readable format to describe the pod
func (c *Client) DescribePod(pod types.PodDiscovery, options DescribeOptions) (string, error) {
	d, err := c.getDescriber(pod.Context, schema.GroupKind{Group: corev1.GroupName, Kind: "Pod"})
	if err != nil {
		return "", err
	}

	return d.Describe(pod.Namespace, pod.Name, describe.DescriberSettings{ShowEvents: options.ShowEvents})
}

// DescribeCronJob returns a human readable format to describe the cronjob
func (c *Client) DescribeCronJob(cronjob types.CronJobDiscovery, options DescribeOptions) (string, error) {
	d, err := c.getDescriber(cronjob.Context, schema.GroupKind{Group: batchv1.GroupName, Kind: "CronJob"})
	if err != nil {
		return "", err
	}

	s, err := d.Describe(cronjob.Namespace, cronjob.Name, describe.DescriberSettings{ShowEvents: options.ShowEvents})
	return c.describeContextInfo(cronjob.Context) + s, err
}

// DescribeJob returns a human readable format to describe the job
func (c *Client) DescribeJob(job types.JobDiscovery, options DescribeOptions) (string, error) {
	d, err := c.getDescriber(job.Context, schema.GroupKind{Group: batchv1.GroupName, Kind: "Job"})
	if err != nil {
		return "", err
	}

	s, err := d.Describe(job.Namespace, job.Name, describe.DescriberSettings{ShowEvents: options.ShowEvents})
	return c.describeContextInfo(job.Context) + s, err
}

// DescribeConfigMap returns a human readable format to describe the configmap
func (c *Client) DescribeConfigMap(configmap types.ConfigMapDiscovery, options DescribeOptions) (string, error) {
	d, err := c.getDescriber(configmap.Context, schema.GroupKind{Group: corev1.GroupName, Kind: "ConfigMap"})
	if err != nil {
		return "", err
	}

	s, err := d.Describe(configmap.Namespace, configmap.Name, describe.DescriberSettings{ShowEvents: options.ShowEvents})
	return c.describeContextInfo(configmap.Context) + s, err
}

// DescribeDeployment returns a human readable format to describe the deployment
func (c *Client) DescribeDeployment(deployment types.DeploymentDiscovery, options DescribeOptions) (string, error) {
	d, err := c.getDescriber(deployment.Context, schema.GroupKind{Group: appsv1.GroupName, Kind: "Deployment"})
	if err != nil {
		return "", err
	}

	s, err := d.Describe(deployment.Namespace, deployment.Name, describe.DescriberSettings{ShowEvents: options.ShowEvents})
	return c.describeContextInfo(deployment.Context) + s, err
}

// DescribeReplicaSet returns a human readable format to describe the replicaset
func (c *Client) DescribeReplicaSet(replicaset types.ReplicaSetDiscovery, options DescribeOptions) (string, error) {
	d, err := c.getDescriber(replicaset.Context, schema.GroupKind{Group: appsv1.GroupName, Kind: "ReplicaSet"})
	if err != nil {
		return "", err
	}

	s, err := d.Describe(replicaset.Namespace, replicaset.Name, describe.DescriberSettings{ShowEvents: options.ShowEvents})
	return c.describeContextInfo(replicaset.Context) + s, err
}
