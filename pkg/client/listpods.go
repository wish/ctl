package client

import (
	"sync"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ContextLogic/ctl/pkg/client/helper"
)

func (c *Client) ListPods(context string, namespace string, options ListOptions) ([]v1.Pod, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	pods, err := cs.CoreV1().Pods(namespace).List(metav1.ListOptions{Limit: options.Limit})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
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

			list, err := c.ListPods(ctx, namespace, options)
			if err != nil { return }

			mutex.Lock()
			for _, x := range list {
				ret = append(ret, PodDiscovery{ctx, x.Namespace, x})
			}
			mutex.Unlock()
		}(ctx)
	}

	wait.Wait()
	return ret, nil
}
