package client

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client/helper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sync"
)

func (c *Client) ListPods(context string, namespace string, options ListOptions) ([]PodDiscovery, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	pods, err := cs.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	items := make([]PodDiscovery, len(pods.Items))
	for i, p := range pods.Items {
		items[i] = PodDiscovery{context, p}
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

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			pods, err := c.ListPods(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, `Could not connect to cluster "%s": %v\n`, ctx, err)
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
	return ret, nil
}
