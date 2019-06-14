package client

import (
	"github.com/ContextLogic/ctl/pkg/client/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetRun(context, namespace string, name string, options GetOptions) (*RunDiscovery, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	job, err := cs.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &RunDiscovery{context, *job}, nil
}

func (c *Client) FindRuns(contexts []string, namespace string, names []string, options ListOptions) ([]RunDiscovery, error) {
	if len(contexts) == 0 {
		contexts = helper.GetContexts()
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	all, err := c.ListRunsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []RunDiscovery

	for _, j := range all {
		if _, ok := positive[j.Name]; ok {
			ret = append(ret, j)
		}
	}

	return ret, nil
}
