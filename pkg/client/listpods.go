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

// ListPods returns a list of all pod that match the query
func (c *Client) ListPods(context string, namespace string, options ListOptions) ([]types.PodDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	pods, err := cs.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.PodDiscovery
	for _, pod := range pods.Items {
		p := types.PodDiscovery{context, pod}
		c.Transform(&p)
		if filter.MatchLabel(p, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(p.Name)) {
			items = append(items, p)
		}
	}
	return items, nil
}

// ListPodsOverContexts is like ListPods but operates over multiple clusters
func (c *Client) ListPodsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.PodDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.PodDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			pods, err := c.ListPods(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				mutex.Lock()
				failed = append(failed, ctx)
				mutex.Unlock()
				return
			}

			mutex.Lock()
			for _, pod := range pods {
				ret = append(ret, pod)
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

// ListPodsOfJob returns a list of all pods belonging to a job
func (c *Client) ListPodsOfJob(contexts []string, namespace, jobName string, options ListOptions) ([]types.PodDiscovery, error) {
	pods, err := c.ListPodsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	job, err := c.findJob(contexts, namespace, jobName, options)
	if err != nil {
		return nil, err
	}

	var ret []types.PodDiscovery
	for _, p := range pods {
		// Check if has owner reference
		for _, o := range p.OwnerReferences {
			if o.UID == job.UID { // matches
				ret = append(ret, p)
				break
			}
		}
	}

	return ret, nil
}
