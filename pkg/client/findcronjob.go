package client

import (
	"context"
	"errors"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helpers for finding a specific cron job
func (c *Client) findCronJob(contexts []string, namespace, name string, options ListOptions) (*types.CronJobDiscovery, error) {
	list, err := c.ListCronJobsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var cron types.CronJobDiscovery
	for _, cj := range list {
		if cj.Name == name {
			cron = cj
			break
		}
	}

	if cron.Name != name { // Pod not found
		return nil, errors.New("cron job not found") // TODO return value
	}

	return &cron, nil
}

func (c *Client) findExactCronJob(contextStr, namespace, name string) (*types.CronJobDiscovery, error) {
	cl, err := c.getContextInterface(contextStr)
	if err != nil {
		return nil, err
	}

	cj, err := cl.BatchV1beta1().CronJobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &types.CronJobDiscovery{contextStr, *cj}, nil
}
