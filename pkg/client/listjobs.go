package client

import (
	"errors"
	"fmt"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"sync"
)

// ListJobs returns a list of all jobs that match the query
func (c *Client) ListJobs(context string, namespace string, options ListOptions) ([]types.JobDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	jobs, err := cs.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.JobDiscovery
	for _, job := range jobs.Items {
		r := types.JobDiscovery{context, job}
		c.extension.Transform(&r)
		if filter.MatchLabel(r, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(r.Name)) {
			items = append(items, r)
		}
	}
	return items, nil
}

// ListJobsOverContexts is like ListJobs but operates over multiple clusters
func (c *Client) ListJobsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.JobDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.extension.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.extension.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.JobDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			jobs, err := c.ListJobs(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				failed = append(failed, ctx)
				return
			}

			mutex.Lock()
			for _, job := range jobs {
				ret = append(ret, job)
			}
			mutex.Unlock()
		}(ctx)
	}

	wait.Wait()
	if failed != nil {
		return ret, errors.New("failed connecting to the following contexts: " + strings.Join(failed, ", "))
	}
	sortObjs(ret)
	return ret, nil
}

// ListJobsOfCronJob returns a list of all jobs belonging to a cron job
func (c *Client) ListJobsOfCronJob(contexts []string, namespace, cronjobName string, options ListOptions) ([]types.JobDiscovery, error) {
	cronjob, err := c.findCronJob(contexts, namespace, cronjobName, options)
	if err != nil {
		return nil, err
	}

	// Assuming that jobs started are in the same location
	list, err := c.ListJobs(cronjob.Context, cronjob.Namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []types.JobDiscovery
	for _, j := range list {
		// Check if has owner reference
		for _, o := range j.OwnerReferences {
			if o.UID == cronjob.UID { // matches
				ret = append(ret, j)
				break
			}
		}
	}

	return ret, nil
}
