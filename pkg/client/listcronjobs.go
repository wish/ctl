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

// ListCronJobs returns a list of all cron jobs that match the query
func (c *Client) ListCronJobs(context string, namespace string, options ListOptions) ([]types.CronJobDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	cronjobs, err := cs.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.CronJobDiscovery
	for _, cj := range cronjobs.Items {
		cjd := types.CronJobDiscovery{context, cj}
		c.Transform(&cjd)
		if filter.MatchLabel(cjd, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(cjd.Name)) { // TODO: Modularize to allow adding more search parameters
			items = append(items, cjd)
		}
	}
	return items, nil
}

// ListCronJobsOverContexts is like ListCronJobs but operates over multiple clusters
func (c *Client) ListCronJobsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.CronJobDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.CronJobDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			list, err := c.ListCronJobs(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				failed = append(failed, ctx)
				return
			}

			mutex.Lock()
			for _, j := range list {
				ret = append(ret, j)
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
