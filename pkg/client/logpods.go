package client

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// Retrieves logs of a single pod (uses first found if multiple)
func (c *Client) LogPodOverContexts(contexts []string, namespace, name, container string, options LogOptions) (*rest.Result, error) {
	pod, container, err := c.findPodWithContainer(contexts, namespace, name, container)
	if err != nil {
		return nil, err
	}

	cl, err := c.getContextInterface(pod.Context)
	if err != nil {
		panic(err.Error())
	}

	req := cl.CoreV1().Pods(pod.Namespace).GetLogs(name, &v1.PodLogOptions{Container: container})
	res := req.Do()
	return &res, nil
}

// Only logs first container if container not specified
// TODO: The usage of this function is odd (support all containers???)
func (c *Client) LogPod(context, namespace, name, container string, options LogOptions) (*rest.Result, error) {
	cl, err := c.getContextInterface(context)
	if err != nil {
		panic(err.Error())
	}

	// Find first container
	if container == "" || namespace == "" {
		pod, err := c.findPod([]string{context}, namespace, name)
		if err != nil {
			return nil, err
		}
		if container == "" {
			container = pod.Spec.Containers[0].Name
		}
		namespace = pod.Namespace
	}

	req := cl.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{Container: container})
	res := req.Do()
	return &res, nil
}
