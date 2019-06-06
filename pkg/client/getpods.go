package client

import (
  // "fmt"
	"sync"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "github.com/ContextLogic/ctl/pkg/client/helper"
)

func (c *Client) GetPod(context, namespace string, name string, options GetOptions) (*v1.Pod, error) {
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil, err
	}
	pod, err := cs.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

// If contexts and namespaces are left blank, then searches through all
func (c *Client) GetPodOverContext(contexts, namespaces []string, name string, options GetOptions) ([]PodDiscovery, error) {
  if len(contexts) == 0 {
    contexts = helper.GetContexts()
  }

	var waitc sync.WaitGroup
	waitc.Add(len(contexts))

	var mutex sync.Mutex // lock for ret
	var ret []PodDiscovery

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

					pod, err := c.GetPod(ctx, ns, name, options)
					if err != nil { return }

					mutex.Lock()
          ret = append(ret, PodDiscovery{ctx, ns, *pod})
					mutex.Unlock()
				}(ns)
			}

			waitn.Wait()
		}(ctx)
	}

	waitc.Wait()
	return ret, nil
}
