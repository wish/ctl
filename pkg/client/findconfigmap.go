package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/types"
)

// Helpers for finding a specific config map
func (c *Client) findConfigMap(contexts []string, namespace, name string, options ListOptions) (*types.ConfigMapDiscovery, error) {
	list, err := c.ListConfigMapsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var cm types.ConfigMapDiscovery
	for _, cj := range list {
		if cj.Name == name {
			cm = cj
			break
		}
	}

	if cm.Name != name { // Pod not found
		return nil, errors.New("configmap not found") // TODO return value
	}

	return &cm, nil
}
