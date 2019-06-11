package client

import (
	"sync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ContextLogic/ctl/pkg/client/helper"
)

func (c *Client) ListPods(context string, namespace string, options ListOptions) ([]PodDiscovery, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	pods, err := cs.CoreV1().Pods(namespace).List(metav1.ListOptions{Limit: options.Limit})
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
			if err != nil { return }

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
