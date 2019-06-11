package client

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

// Retrieves logs of a single pod (uses first found if multiple)
func (c *Client) LogPod(contexts []string, namespace, name, container string, options LogOptions) (*rest.Result, error) {
	pod, container, err := c.findPodWithContainer(contexts, namespace, name, container)
	if err != nil {
		return nil, err
	}

	cl, err := c.getContextClientset(pod.Context)
	if err != nil {
		panic(err.Error())
	}

	req := cl.CoreV1().Pods(pod.Namespace).GetLogs(name, &v1.PodLogOptions{Container: container})
	res := req.Do()
	return &res, nil
}
