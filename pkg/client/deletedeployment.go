package client

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteDeployment deletes a single deployment. Requires exact location.
func (c *Client) DeleteDeployment(context, namespace, name string, options DeleteOptions) error {
	cl, err := c.getContextInterface(context)
	if err != nil {
		return err
	}

	var deleteOptions *metav1.DeleteOptions
	if options.Now {
		var one int64 = 1
		deleteOptions = &metav1.DeleteOptions{GracePeriodSeconds: &one}
	}

	return cl.ExtensionsV1beta1().Deployments(namespace).Delete(name, deleteOptions)
}
