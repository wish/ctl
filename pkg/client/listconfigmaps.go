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

// ListConfigMaps returns a list of all configmap that match the query
func (c *Client) ListConfigMaps(context string, namespace string, options ListOptions) ([]types.ConfigMapDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	configmaps, err := cs.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.ConfigMapDiscovery
	for _, configmap := range configmaps.Items {
		cm := types.ConfigMapDiscovery{context, configmap}
		c.Transform(&cm)
		if filter.MatchLabel(cm, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(cm.Name)) {
			items = append(items, cm)
		}
	}
	return items, nil
}

// ListConfigMapsOverContexts is like ListConfigMaps but operates over multiple clusters
func (c *Client) ListConfigMapsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.ConfigMapDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.ConfigMapDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			configmaps, err := c.ListConfigMaps(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				mutex.Lock()
				failed = append(failed, ctx)
				mutex.Unlock()
				return
			}

			mutex.Lock()
			for _, configmap := range configmaps {
				ret = append(ret, configmap)
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
