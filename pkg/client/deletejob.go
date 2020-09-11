package client

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteJob deletes a single job. Requires exact location.
func (c *Client) DeleteJob(contextStr, namespace, name string, options DeleteOptions) error {
	cl, err := c.getContextInterface(contextStr)
	if err != nil {
		return err
	}

	deleteOptions := metav1.DeleteOptions{}
	if options.Now {
		var one int64 = 1
		deleteOptions.GracePeriodSeconds = &one
	}
	if options.DeletionPropagation {
		bg := metav1.DeletePropagationBackground
		deleteOptions.PropagationPolicy = &bg
	}

	return cl.BatchV1().Jobs(namespace).Delete(context.TODO(), name, deleteOptions)
}
