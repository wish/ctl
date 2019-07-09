package client

import (
	"errors"
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"sync"
)

func (c *Client) ListPods(context string, namespace string, options ListOptions) ([]PodDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	pods, err := cs.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]PodDiscovery, len(pods.Items))
	for i, pod := range pods.Items {
		items[i] = PodDiscovery{context, pod}
	}
	return items, nil
}

func (c *Client) ListPodsOverContexts(contexts []string, namespace string, options ListOptions) ([]PodDiscovery, error) {
	if len(contexts) == 0 {
		contexts = helper.GetContexts()
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []PodDiscovery
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			pods, err := c.ListPods(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				failed = append(failed, ctx)
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
		return ret, errors.New("Failed connecting to the following contexts: " + strings.Join(failed, ", "))
	}
	return ret, nil
}

func (c *Client) ListPodsOfRun(contexts []string, namespace, runName string, options ListOptions) ([]PodDiscovery, error) {
	pods, err := c.ListPodsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	run, err := c.findRun(contexts, namespace, runName)
	if err != nil {
		return nil, err
	}

	var ret []PodDiscovery
	for _, p := range pods {
		// Check if has owner reference
		for _, o := range p.OwnerReferences {
			if o.UID == run.UID { // matches
				ret = append(ret, p)
				break
			}
		}
	}

	return ret, nil
}
