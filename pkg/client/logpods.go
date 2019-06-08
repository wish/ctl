package client

import (
  "errors"
  "fmt"
  "k8s.io/api/core/v1"
)

// https://stackoverflow.com/questions/32983228/kubernetes-go-client-api-for-log-of-a-particular-pod
func (c *Client) LogPod(contexts []string, namespace, name, container string, options LogOptions) error {
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
    return errors.New("Pod not found") // TODO return value
  }

  // Container
  if container == "" { // No container specified
    if len(pod.Spec.Containers) == 1 {
      container = pod.Spec.Containers[0].Name
    } else {
      for _, c := range pod.Spec.Containers {
        fmt.Println(c.Name)
      }
      container = pod.Spec.Containers[0].Name
      // return nil, errors.New("There are multiple containers and none was specified!")
    }
  }
  cl, err := c.getContextClientset(pod.Context)
  if err != nil {
    panic(err.Error())
  }
  req := cl.CoreV1().Pods(pod.Namespace).GetLogs(name, &v1.PodLogOptions{Container: container})
  raw, err := req.DoRaw()
  if err != nil {
    panic(err.Error())
  }
  fmt.Println(string(raw))

  return nil
}
