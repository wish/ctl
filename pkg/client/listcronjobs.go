package client

import (
	"sync"
	// "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ContextLogic/ctl/pkg/client/helper"
)

func (c *Client) ListCronJobs(context string, namespace string, options ListOptions) ([]CronJobDiscovery, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	cronjobs, err := cs.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{Limit: options.Limit})
	if err != nil {
		return nil, err
	}
	items := make([]CronJobDiscovery, len(cronjobs.Items))
	for i, j := range cronjobs.Items {
		items[i] = CronJobDiscovery{context, j}
	}
	return items, nil
}

func (c *Client) ListCronJobsOverContexts(contexts []string, namespace string, options ListOptions) ([]CronJobDiscovery, error) {
	if len(contexts) == 0 {
		contexts = helper.GetContexts()
	}

	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []CronJobDiscovery

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			list, err := c.ListCronJobs(ctx, namespace, options)
			if err != nil { return }

			mutex.Lock()
			for _, j := range list {
				ret = append(ret, j)
			}
			mutex.Unlock()
		}(ctx)
	}

	wait.Wait()
	return ret, nil
}
