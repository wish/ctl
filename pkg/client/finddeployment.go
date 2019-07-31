package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/types"
)

// Helpers for finding a specific config map
func (c *Client) findDeployment(contexts []string, namespace, name string, options ListOptions) (*types.DeploymentDiscovery, error) {
	list, err := c.ListDeploymentsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var d types.DeploymentDiscovery
	for _, i := range list {
		if i.Name == name {
			d = i
			break
		}
	}

	if d.Name != name { // Pod not found
		return nil, errors.New("deployment not found") // TODO return value
	}

	return &d, nil
}
