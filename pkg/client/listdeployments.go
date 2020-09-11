package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"sync"
)

// ListDeployments returns a list of all deployments that match the query
func (c *Client) ListDeployments(contextStr string, namespace string, options ListOptions) ([]types.DeploymentDiscovery, error) {
	cs, err := c.getContextInterface(contextStr)
	if err != nil {
		return nil, err
	}
	deployments, err := cs.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var items []types.DeploymentDiscovery
	for _, deployment := range deployments.Items {
		cm := types.DeploymentDiscovery{contextStr, deployment}
		c.Transform(&cm)
		if filter.MatchLabel(cm, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(cm.Name)) {
			items = append(items, cm)
		}
	}
	return items, nil
}

// ListDeploymentsOverContexts is like ListDeployments but operates over multiple clusters
func (c *Client) ListDeploymentsOverContexts(contexts []string, namespace string, options ListOptions) ([]types.DeploymentDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetFilteredContexts(options.LabelMatch)
	} else {
		contexts = c.FilterContexts(contexts, options.LabelMatch)
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []types.DeploymentDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			deployments, err := c.ListDeployments(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				mutex.Lock()
				failed = append(failed, ctx)
				mutex.Unlock()
				return
			}

			mutex.Lock()
			for _, deployment := range deployments {
				ret = append(ret, deployment)
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
