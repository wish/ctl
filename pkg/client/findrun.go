package client

import (
	"errors"
	"github.com/ContextLogic/ctl/pkg/client/types"
)

// Helpers for finding a specific run
func (c *Client) findRun(contexts []string, namespace, name string) (*types.RunDiscovery, error) {
	list, err := c.ListRunsOverContexts(contexts, namespace, ListOptions{})
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
		return nil, errors.New("Cron job run with name \"" + name + "\" not found") // TODO return value
	}

	return &run, nil
}
