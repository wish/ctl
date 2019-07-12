package client

import "github.com/ContextLogic/ctl/pkg/client/filter"

type ListOptions struct {
	filter.LabelMatch
}

type GetOptions struct {
	filter.LabelMatch
}

type LogOptions struct {
	filter.LabelMatch
}
