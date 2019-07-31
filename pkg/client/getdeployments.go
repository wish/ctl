package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeployment returns a single deployment
func (c *Client) GetDeployment(context, namespace string, name string, options GetOptions) (*types.DeploymentDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	// REVIEW: In the future it will be useful to have a function to convert client.GetOptions -> metav1.GetOptions
	configmap, err := cs.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	d := types.DeploymentDiscovery{context, *configmap}
	c.extension.Transform(&d)
	if !filter.MatchLabel(d, options.LabelMatch) {
		return nil, errors.New("found object does not satisfy filters")
	}
	return &d, nil
}

// FindDeployments simultaneously searches for multiple deployments and returns all results
func (c *Client) FindDeployments(contexts []string, namespace string, names []string, options ListOptions) ([]types.DeploymentDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.extension.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.extension.FilterContexts(contexts, options.LabelMatch)
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	all, err := c.ListDeploymentsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []types.DeploymentDiscovery

	for _, d := range all {
		if _, ok := positive[d.Name]; ok {
			ret = append(ret, d)
		}
	}

	return ret, nil
}
