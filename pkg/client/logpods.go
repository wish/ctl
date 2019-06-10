package client

import (
  "errors"
  "k8s.io/api/core/v1"
  "strings"
  "k8s.io/client-go/rest"
)

func (c *Client) LogPod(contexts []string, namespace, name, container string, options LogOptions) (*rest.Result, error) {
  // Find pod
  // REVIEW: change this to use get instead of list
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

  // Container
  if container == "" { // No container specified
    if len(pod.Spec.Containers) == 1 {
      container = pod.Spec.Containers[0].Name
    } else {
      conts := make([]string, len(pod.Spec.Containers))
      for i, c := range pod.Spec.Containers {
        conts[i] = c.Name
      }
      return nil, errors.New("No container was specified! Choose one of the containers: " + strings.Join(conts, ", "))
    }
  }

  cl, err := c.getContextClientset(pod.Context)
  if err != nil {
    panic(err.Error())
  }

  req := cl.CoreV1().Pods(pod.Namespace).GetLogs(name, &v1.PodLogOptions{Container: container})
  res := req.Do()
  return &res, nil
}
