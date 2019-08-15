package client

import (
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeletePod deletes a single pod. Requires exact location.
func (c *Client) DeletePod(context, namespace, name string, options DeleteOptions) error {
  cl, err := c.getContextInterface(context)
  if err != nil {
    return err
  }

  var deleteOptions *metav1.DeleteOptions
  if options.Now {
    var one int64 = 1
    deleteOptions = &metav1.DeleteOptions{GracePeriodSeconds: &one}
  }
  return cl.CoreV1().Pods(namespace).Delete(name, deleteOptions)
}
