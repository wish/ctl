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

// ListReplicaSets returns a list of all replicaset that match the query
func (c *Client) ListReplicaSets(context string, namespace string, options ListOptions) ([]types.ReplicaSetDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	replicasets, err := cs.ExtensionsV1beta1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.ReplicaSetDiscovery
	for _, replicaset := range replicasets.Items {
		cm := types.ReplicaSetDiscovery{context, replicaset}
		c.Transform(&cm)
		if filter.MatchLabel(cm, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(cm.Name)) {
			items = append(items, cm)
		}
	}
	return items, nil
}

// ListReplicaSetsOverContexts is like ListReplicaSets but operates over multiple clusters
func (c *Client) ListReplicaSetsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.ReplicaSetDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.ReplicaSetDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			replicasets, err := c.ListReplicaSets(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				mutex.Lock()
				failed = append(failed, ctx)
				mutex.Unlock()
				return
			}

			mutex.Lock()
			for _, replicaset := range replicasets {
				ret = append(ret, replicaset)
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
