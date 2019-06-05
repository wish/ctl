package kron

import (
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) List(options ListOptions) ([]v1beta1.CronJob, error) {
	cronjobs, err := c.clientset.BatchV1beta1().CronJobs("").List(metav1.ListOptions{Limit: options.Limit})
	if err != nil {
		return nil, err
	}

	// Add more search options.
	// Search my keyword will have to be done after querying.
	// Thus, limit will have to be changed post-processing

	return cronjobs.Items, nil
}
