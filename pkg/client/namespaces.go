// REVIEW: This could belong in helper.go or client.go instead of its own file
package client

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Helper for getting all namespaces
// REVIEW: Consider removing as this method is currently not being used
func (c *Client) GetNamespaces(context string) []string {
	// Default options
	cs, err := c.getContextClientset(context)
	if err != nil {
		return nil
	}
	namespaces, err := cs.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		// TODO: Do some logging on why namespaces couldn't be found
		fmt.Println(err.Error())
		return []string{}
	}

	list := make([]string, len(namespaces.Items))

	for _, n := range namespaces.Items { // Currently ignoring mappings
		list = append(list, n.Name)
	}

	return list
}
