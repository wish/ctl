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

func (c *Client) ListCronJobs(context string, namespace string, options ListOptions) ([]CronJobDiscovery, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	cronjobs, err := cs.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
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
	var failed []string

	for _, ctx := range contexts {
		go func(ctx string) {
			defer wait.Done()

			list, err := c.ListCronJobs(ctx, namespace, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not connect to cluster \"%s\": %v\n", ctx, err)
				failed = append(failed, ctx)
				return
			}

			mutex.Lock()
			for _, j := range list {
				ret = append(ret, j)
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
