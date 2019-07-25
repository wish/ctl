package client

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// LogPodOverContexts retrieves logs of a single pod (uses first found if multiple)
func (c *Client) LogPodOverContexts(contexts []string, namespace, name, container string, options LogOptions) (*rest.Request, error) {
	pod, container, err := c.findPodWithContainer(contexts, namespace, name, container, ListOptions{options.LabelMatch, nil})
	if err != nil {
		return nil, err
	}

	cl, err := c.getContextInterface(pod.Context)
	if err != nil {
		panic(err.Error())
	}

	req := cl.CoreV1().Pods(pod.Namespace).GetLogs(name, &v1.PodLogOptions{Container: container, Follow: options.Follow})
	return req, nil
}

// LogPod retrieves logs from a container of a pod.
// Operates on the first container if none specified.
// TODO: The usage of this function is odd (support all containers???)
func (c *Client) LogPod(context, namespace, name, container string, options LogOptions) (*rest.Request, error) {
	cl, err := c.getContextInterface(context)
	if err != nil {
		panic(err.Error())
	}

	// Find first container
	if container == "" || namespace == "" {
		pod, err := c.findPod([]string{context}, namespace, name, ListOptions{options.LabelMatch, nil})
		if err != nil {
			return nil, err
		}
		if container == "" {
			container = pod.Spec.Containers[0].Name
		}
		namespace = pod.Namespace
	}

	req := cl.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{Container: container, Follow: options.Follow})
	return req, nil
}
