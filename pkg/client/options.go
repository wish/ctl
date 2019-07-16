package client

import "github.com/wish/ctl/pkg/client/filter"

type ListOptions struct {
	filter.LabelMatch
}

type GetOptions struct {
	filter.LabelMatch
}

type LogOptions struct {
	filter.LabelMatch
}
