package client

import "github.com/wish/ctl/pkg/client/filter"

// Currently, all three options are the same and only support filtering.

// ListOptions is used to specific filtering on list operations
type ListOptions struct {
	filter.LabelMatch
}

// GetOptions is used to specific filtering on get operations
type GetOptions struct {
	filter.LabelMatch
}

// LogOptions is used to specific filtering on log operations
type LogOptions struct {
	filter.LabelMatch
}
