package client

import (
	"sync"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) ListCronJobs(context string, options ListOptions) ([]v1beta1.CronJob, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	cronjobs, err := cs.BatchV1beta1().CronJobs("").List(metav1.ListOptions{Limit: options.Limit})
	if err != nil {
		return nil, err
	}

	// Add more search options.
	// Search my keyword will have to be done after querying.
	// Thus, limit will have to be changed post-processing

	return cronjobs.Items, nil
}

func (c *Client) ListCronJobsOverContexts(contexts []string, options ListOptions) ([]CronJobDiscovery, error) {
	var wait sync.WaitGroup
	wait.Add(len(contexts))

	var mutex sync.Mutex
	var ret []CronJobDiscovery

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			list, err := c.ListCronJobs(ctx, options)
			if err != nil { return }

			mutex.Lock()
			for _, x := range list {
				ret = append(ret, CronJobDiscovery{ctx, x.Namespace, x})
			}
			mutex.Unlock()
		}(ctx)
	}

	wait.Wait()
	return ret, nil
}
