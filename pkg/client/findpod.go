package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wish/ctl/pkg/client/types"
)

// Helpers for finding a specific pod
func (c *Client) findPod(contexts []string, namespace, name string, options ListOptions) (*types.PodDiscovery, error) {
	list, err := c.ListPodsOverContexts(contexts, namespace, options)
	if err != nil {
		return nil, err
	}

	var pod types.PodDiscovery
	for _, p := range list {
		if p.Name == name {
			pod = p
			break
		}
	}

	if pod.Name != name { // Pod not found
		return nil, errors.New("pod not found") // TODO return value
	}

	return &pod, nil
}

func (c *Client) FindPodWithContainer(contexts []string, namespace, name, optionalContainer string, options ListOptions) (pod *types.PodDiscovery, container string, err error) {
	pod, err = c.findPod(contexts, namespace, name, options)
	if err != nil {
		return
	}

	if optionalContainer == "" {
		if len(pod.Spec.Containers) > 0 {
			container = pod.Spec.Containers[0].Name
			var s []string
			for _, c := range pod.Spec.Containers {
				s = append(s, c.Name)
			}
			fmt.Println("Available containers are:", strings.Join(s, ", "))
			fmt.Println("No container specified, defaulting to the first container:", container)
		} else {
			err = errors.New("there are no containers on this pod")
		}
	} else {
		container = optionalContainer
	}

	return
}
