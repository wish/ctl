package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/types"
)

// Helpers for finding a specific run
func (c *Client) findRun(contexts []string, namespace, name string, options ListOptions) (*types.RunDiscovery, error) {
	list, err := c.ListRunsOverContexts(contexts, namespace, options)
	if err != nil {
		panic(err.Error())
	}

	var run types.RunDiscovery
	for _, r := range list {
		if r.Name == name {
			run = r
			break
		}
	}

	if run.Name != name { // Pod not found
		return nil, errors.New("cron job run with name \"" + name + "\" not found") // TODO return value
	}

	return &run, nil
}
