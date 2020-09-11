package client

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteReplicaSet deletes a single replicaset. Requires exact location.
func (c *Client) DeleteReplicaSet(contextStr, namespace, name string, options DeleteOptions) error {
	cl, err := c.getContextInterface(contextStr)
	if err != nil {
		return err
	}

	var deleteOptions metav1.DeleteOptions
	if options.Now {
		var one int64 = 1
		deleteOptions = metav1.DeleteOptions{GracePeriodSeconds: &one}
	}

	return cl.AppsV1().ReplicaSets(namespace).Delete(context.TODO(), name, deleteOptions)
}
