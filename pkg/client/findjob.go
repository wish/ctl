package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/types"
)

// Helper for finding a specific job
func (c *Client) findJob(contexts []string, namespace, name string, options ListOptions) (*types.JobDiscovery, error) {
	list, err := c.ListJobsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var job types.JobDiscovery
	for _, j := range list {
		if j.Name == name {
			job = j
			break
		}
	}

	if job.Name != name { // Pod not found
		return nil, errors.New("job with name \"" + name + "\" not found") // TODO return value
	}

	return &job, nil
}
