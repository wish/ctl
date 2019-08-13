package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetReplicaSet returns a single replicaset
func (c *Client) GetReplicaSet(context, namespace string, name string, options GetOptions) (*types.ReplicaSetDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	// REVIEW: In the future it will be useful to have a function to convert client.GetOptions -> metav1.GetOptions
	replicaset, err := cs.ExtensionsV1beta1().ReplicaSets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	r := types.ReplicaSetDiscovery{context, *replicaset}
	c.Transform(&r)
	if !filter.MatchLabel(r, options.LabelMatch) {
		return nil, errors.New("found object does not satisfy filters")
	}
	return &r, nil
}

// FindReplicaSets simultaneously searches for multiple configmaps and returns all results
func (c *Client) FindReplicaSets(contexts []string, namespace string, names []string, options ListOptions) ([]types.ReplicaSetDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.FilterContexts(contexts, options.LabelMatch)
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	all, err := c.ListReplicaSetsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []types.ReplicaSetDiscovery

	for _, j := range all {
		if _, ok := positive[j.Name]; ok {
			ret = append(ret, j)
		}
	}

	return ret, nil
}
