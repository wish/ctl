package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/types"
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
		return nil, errors.New("Cron job not found") // TODO return value
	}

	return &cron, nil
}
