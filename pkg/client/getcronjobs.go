package client

import (
	"errors"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCronJob returns a single cron job
func (c *Client) GetCronJob(context, namespace string, name string, options GetOptions) (*types.CronJobDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	// REVIEW: In the future it will be useful to have a function to convert client.GetOptions -> metav1.GetOptions
	cronjob, err := cs.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	d := types.CronJobDiscovery{context, *cronjob}
	if !filter.MatchLabel(d, options.LabelMatch) {
		return nil, errors.New("found object does not satisfy filters")
	}
	c.forger.Transform(&d)
	return &d, nil
}

// FindCronJobs simultaneously searches for multiple cron jobs and returns all results
func (c *Client) FindCronJobs(contexts []string, namespace string, names []string, options ListOptions) ([]types.CronJobDiscovery, error) {
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

	var ret []types.CronJobDiscovery

	for _, j := range all {
		if _, ok := positive[j.Name]; ok {
			ret = append(ret, j)
			c.forger.Transform(&ret[len(ret)-1])
		}
	}

	return ret, nil
}
