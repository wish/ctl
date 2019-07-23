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

// ListRuns returns a list of all jobs that match the query
func (c *Client) ListRuns(context string, namespace string, options ListOptions) ([]types.RunDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	runs, err := cs.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.RunDiscovery
	for _, run := range runs.Items {
		r := types.RunDiscovery{context, run}
		c.extension.Transform(&r)
		if filter.MatchLabel(r, options.LabelMatch) {
			items = append(items, r)
		}
	}
	return items, nil
}

// ListRunsOverContexts is like ListRuns but operates over multiple clusters
func (c *Client) ListRunsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.RunDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.extension.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.extension.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.RunDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			runs, err := c.ListRuns(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				failed = append(failed, ctx)
				return
			}

			mutex.Lock()
			for _, run := range runs {
				ret = append(ret, run)
			}
			mutex.Unlock()
		}(ctx)
	}

	wait.Wait()
	if failed != nil {
		return ret, errors.New("failed connecting to the following contexts: " + strings.Join(failed, ", "))
	}
	return ret, nil
}

// ListRunsOfCronJob returns a list of all jobs belonging to a cron job
func (c *Client) ListRunsOfCronJob(contexts []string, namespace, cronjobName string, options ListOptions) ([]types.RunDiscovery, error) {
	cronjob, err := c.findCronJob(contexts, namespace, cronjobName, options)
	if err != nil {
		return nil, err
	}

	// Assuming that jobs started are in the same location
	list, err := c.ListRuns(cronjob.Context, cronjob.Namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []types.RunDiscovery
	for _, r := range list {
		// Check if has owner reference
		for _, o := range r.OwnerReferences {
			if o.UID == cronjob.UID { // matches
				ret = append(ret, r)
				break
			}
		}
	}

	return ret, nil
}
