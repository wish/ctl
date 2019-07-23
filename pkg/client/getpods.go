package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod returns a single pod
func (c *Client) GetPod(context, namespace string, name string, options GetOptions) (*types.PodDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	pod, err := cs.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	d := types.PodDiscovery{context, *pod}
	c.extension.Transform(&d)
	if !filter.MatchLabel(d, options.LabelMatch) {
		return nil, errors.New("found object does not satisfy filters")
	}
	return &d, nil
}

// FindPods simultaneously searches for multiple pods and returns all results
func (c *Client) FindPods(contexts []string, namespace string, names []string, options ListOptions) ([]types.PodDiscovery, error) {
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

	all, err := c.ListPodsOverContexts(contexts, namespace, options)

	var ret []types.PodDiscovery

	for _, p := range all {
		if _, ok := positive[p.Name]; ok {
			ret = append(ret, p)
		}
	}

	return ret, err
}
