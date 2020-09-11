package client

import (
	"context"
	"errors"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetJob returns a single job
func (c *Client) GetJob(contextStr, namespace string, name string, options GetOptions) (*types.JobDiscovery, error) {
	cs, err := c.getContextInterface(contextStr)
	if err != nil {
		return nil, err
	}
	job, err := cs.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	d := types.JobDiscovery{contextStr, *job}
	c.Transform(&d)
	if !filter.MatchLabel(d, options.LabelMatch) {
		return nil, errors.New("found object does not satisfy filters")
	}
	return &d, nil
}

// FindJobs simultaneously searches for multiple jobs and returns all results
func (c *Client) FindJobs(contexts []string, namespace string, names []string, options ListOptions) ([]types.JobDiscovery, error) {
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

	all, err := c.ListJobsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []types.JobDiscovery

	for _, j := range all {
		if _, ok := positive[j.Name]; ok {
			ret = append(ret, j)
		}
	}

	return ret, nil
}
