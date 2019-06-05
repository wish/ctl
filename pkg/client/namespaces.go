// REVIEW: This could belong in helper.go or client.go instead of its own file
package client

import (
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper for getting all namespaces
func (c *Client) GetNamespaces() []string {
	// Default options
	namespaces, err := c.clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	list := make([]string, len(namespaces.Items))

	for _, n := range namespaces.Items { // Currently ignoring mappings
		list = append(list, n.Name)
	}

	return list
}
