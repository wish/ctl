package kron

import (
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetOptions currently does not support any functionality
// so Get does not use the parameter
// options is left as a parameter for consistency
// REVIEW: what namespace to search in?
func (c *Client) Get(name string, options GetOptions) (*v1beta1.CronJob, error) {
	cronjob, err := c.clientset.BatchV1beta1().CronJobs("default").Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return cronjob, nil
}
