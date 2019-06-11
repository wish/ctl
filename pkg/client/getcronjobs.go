package client

import (
	// "fmt"
	// "sync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ContextLogic/ctl/pkg/client/helper"
)

// GetOptions currently does not support any functionality
// so Get does not use the parameter
// options is left as a parameter for consistency
// REVIEW: what namespace to search in?
func (c *Client) GetCronJob(context, namespace string, name string, options GetOptions) (*CronJobDiscovery, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	cronjob, err := cs.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &CronJobDiscovery{context, *cronjob}, nil
}

func (c *Client) FindCronJobs(contexts []string, namespace string, names []string, options ListOptions) ([]CronJobDiscovery, error) {
	if len(contexts) == 0 {
		contexts = helper.GetContexts()
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	var ret []CronJobDiscovery

	all, err := c.ListCronJobsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	for _, p := range all {
		if _, ok := positive[p.Name]; ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
