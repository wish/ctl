package client

import (
	"errors"
)

// Helpers for finding a specific run
func (c *Client) findRun(contexts []string, namespace, name string) (*RunDiscovery, error) {
	list, err := c.ListRunsOverContexts(contexts, namespace, ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var run RunDiscovery
	for _, r := range list {
		if r.Name == name {
			run = r
			break
		}
	}

	if run.Name != name { // Pod not found
		return nil, errors.New("Run not found") // TODO return value
	}

	return &run, nil
}
