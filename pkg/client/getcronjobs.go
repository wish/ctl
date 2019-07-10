package client

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetCronJob(context, namespace string, name string, options GetOptions) (*CronJobDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	// REVIEW: In the future it will be useful to have a function to convert client.GetOptions -> metav1.GetOptions
	cronjob, err := cs.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &CronJobDiscovery{context, *cronjob}, nil
}

func (c *Client) FindCronJobs(contexts []string, namespace string, names []string, options ListOptions) ([]CronJobDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetAllContexts()
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	all, err := c.ListCronJobsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []CronJobDiscovery

	for _, p := range all {
		if _, ok := positive[p.Name]; ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
