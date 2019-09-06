package client

import (
	"github.com/wish/ctl/pkg/client/types"
	"github.com/wish/ctl/pkg/client/filter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"errors"
)

func (c *Client) GetCRDs(context string, namespace string, name string, options ListOptions) (*types.CrdDiscovery, error) {

	crdClient, err := c.getCrdClientSet(context)
	if err != nil {
		return nil, err
	}
	crd, err := crdClient.ApiextensionsV1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})

	d := types.CrdDiscovery{context, *crd}
	c.Transform(&d)
	if !filter.MatchLabel(d, options.LabelMatch) {
		return nil, errors.New("found object does not satisfy filters")
	}
	return &d, nil
	//if err != nil {
	//	return nil, err
	//}
	//var items []types.PodDiscovery
	//for _, pod := range pods.Items {
	//	p := types.PodDiscovery{context, pod}
	//	c.Transform(&p)
	//	if filter.MatchLabel(p, options.LabelMatch) && (options.Search == nil || options.Search.MatchString(p.Name)) {
	//		items = append(items, p)
	//	}
	//}
	//return items, nil
}
