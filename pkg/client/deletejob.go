package client

import (
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteJob deletes a single job. Requires exact location.
func (c *Client) DeleteJob(context, namespace, name string, options DeleteOptions) error {
  cl, err := c.getContextInterface(context)
  if err != nil {
    return err
  }

  var deleteOptions *metav1.DeleteOptions
  if options.Now {
    var one int64 = 1
    deleteOptions = &metav1.DeleteOptions{GracePeriodSeconds: &one}
  }

  return cl.BatchV1().Jobs(namespace).Delete(name, deleteOptions)
}
