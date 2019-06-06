package client

import (
	"fmt"
	"sync"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/ContextLogic/ctl/pkg/client/helper"
)

// GetOptions currently does not support any functionality
// so Get does not use the parameter
// options is left as a parameter for consistency
// REVIEW: what namespace to search in?
func (c *Client) GetCronJob(context, namespace string, name string, options GetOptions) (*v1beta1.CronJob, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	cronjob, err := cs.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cronjob, nil
}

func (c *Client) FindCronJobs(contexts, namespaces, names []string, options ListOptions) ([]CronJobDiscovery, error) {
	if len(contexts) == 0 {
		contexts = helper.GetContexts()
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	var ret []CronJobDiscovery

	for _, ctx := range contexts {
		nss := namespaces
		if len(nss) == 0 {
			nss = []string{""}
		}
		for _, ns := range nss {
			cronjobs, err := c.ListCronJobs(ctx, ns, options)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			for _, cronjob := range cronjobs {
				if _, ok := positive[cronjob.Name]; ok {
					ret = append(ret, CronJobDiscovery{ctx, cronjob})
				}
			}
		}
	}

	return ret, nil
}

// DEPRECATED; TODO: Remove this method
// If contexts and namespaces are left blank, then searches through all
func (c *Client) GetCronJobOverMultiple(contexts, namespaces []string, name string, options GetOptions) ([]CronJobDiscovery, error) {
	var waitc sync.WaitGroup
	waitc.Add(len(contexts))

	var mutex sync.Mutex // lock for ret
	var ret []CronJobDiscovery

	for _, ctx := range contexts {
		go func(ctx string) {
			defer waitc.Done()

			nss := namespaces
			if len(nss) == 0 {
				nss = c.GetNamespaces(ctx)
			}

			var waitn sync.WaitGroup
			waitn.Add(len(nss))

			for _, ns := range nss {
				go func(ns string) {
					defer waitn.Done()

					cronjob, err := c.GetCronJob(ctx, ns, name, options)
					if err != nil { return }

					mutex.Lock()
					ret = append(ret, CronJobDiscovery{ctx, *cronjob})
					mutex.Unlock()
				}(ns)
			}

			waitn.Wait()
		}(ctx)
	}

	waitc.Wait()
	return ret, nil
}
