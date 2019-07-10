package client

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetPod(context, namespace string, name string, options GetOptions) (*PodDiscovery, error) {
	cs, err := c.getContextInterface(context)
	if err != nil {
		return nil, err
	}
	pod, err := cs.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &PodDiscovery{context, *pod}, nil
}

func (c *Client) FindPods(contexts []string, namespace string, names []string, options ListOptions) ([]PodDiscovery, error) {
	if len(contexts) == 0 {
		contexts = c.GetAllContexts()
	}
	// Creating set of names
	positive := make(map[string]struct{})
	for _, name := range names {
		positive[name] = struct{}{}
	}

	all, err := c.ListPodsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var ret []PodDiscovery

	for _, p := range all {
		if _, ok := positive[p.Name]; ok {
			ret = append(ret, p)
		}
	}

	return ret, nil
}
