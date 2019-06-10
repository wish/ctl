package client

import (
  "errors"
  "strings"
)

// Helper for finding a specific pod
func (c *Client) findPod(contexts []string, namespace, name string) (*PodDiscovery, error) {
  list, err := c.ListPodsOverContexts(contexts, namespace, ListOptions{})
  if err != nil {
    panic(err.Error())
  }

  var pod PodDiscovery
  for _, p := range list {
    if p.Name == name {
      pod = p
      break
    }
  }

  if pod.Name != name { // Pod not found
    return nil, errors.New("Pod not found") // TODO return value
  }

  return &pod, nil
}

func (c *Client) findPodWithContainer(contexts []string, namespace, name, container string) (*PodDiscovery, string, error) {
  pod, err := c.findPod(contexts, namespace, name)
  if err != nil {
    return pod, "", err
  }

  // Check for container
  if container == "" { // No container specified
    if len(pod.Spec.Containers) == 1 {
      container = pod.Spec.Containers[0].Name
    } else {
      conts := make([]string, len(pod.Spec.Containers))
      for i, c := range pod.Spec.Containers {
        conts[i] = c.Name
      }
      return nil, "", errors.New("No container was specified! Choose one of the containers: " + strings.Join(conts, ", "))
    }
  }

  return pod, container, nil
}
