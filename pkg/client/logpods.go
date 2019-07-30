package client

import (
	"fmt"
	"github.com/wish/ctl/pkg/client/logsync"
	"io"
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

// LogPodsOverContexts retrieves logs of multiple pods (uses first found if multiple)
func (c *Client) LogPodsOverContexts(contexts []string, namespace, container string, options LogOptions) (io.Reader, error) {
	pods, err := c.ListPodsOverContexts(contexts, namespace, ListOptions{options.LabelMatch, options.Search})
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %d pods\n", len(pods))

	readers := make([]io.Reader, len(pods))

	// Choose container
	for i, pod := range pods {
		var req *rest.Request

		cl, _ := c.getContextInterface(pod.Context)
		// detect container
		if container == "" {
			req = cl.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{Container: pod.Spec.Containers[0].Name, Timestamps: true})
		} else {
			req = cl.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{Container: container, Timestamps: true})
		}

		readCloser, err := req.Stream()
		if err != nil {
			return nil, err
		}

		readers[i] = readCloser
	}

	fmt.Printf("Opened %d connections to pods\n", len(readers))
	return logsync.Sync(readers), nil
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
