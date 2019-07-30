package client

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCtlExt returns the ctl extensions for each cluster
// It maps context (cluster name) -> ConfigMap data
// Return data: no entry for couldn't connect, nil for couldn't open, else actual data
func (c *Client) GetCtlExt() map[string]map[string]string {
	m := make(map[string]map[string]string)

	for _, ctx := range c.GetAllContexts() {
		ci, err := c.getContextInterface(ctx)
		if err != nil { // don't error
			continue
		}

		// Check if the cluster is connectable by attempting to fetch all configmaps
		_, err = ci.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
		if err != nil {
			continue
		}

		cf, err := ci.CoreV1().ConfigMaps("kube-system").Get("ctl-config", metav1.GetOptions{})
		if err != nil {
			m[ctx] = nil
			continue
		}
		m[ctx] = cf.Data
	}
	return m
}
