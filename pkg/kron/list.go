package kron

import (
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/api/batch/v1beta1"
)

func (c *Client) List() ([]v1beta1.CronJob, error) {
  cronjobs, err := c.clientset.BatchV1beta1().CronJobs("").List(metav1.ListOptions{})
  if err != nil {
    return nil, err
  }
  return cronjobs.Items, nil
}
