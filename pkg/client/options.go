package client

import (
	"github.com/wish/ctl/pkg/client/filter"
	"regexp"
)

// ListOptions is used to specific filtering on list operations
type ListOptions struct {
	filter.LabelMatch
	Search *regexp.Regexp
}

// GetOptions is used to specific filtering on get operations
type GetOptions struct {
	filter.LabelMatch
}

// LogOptions is used to specific filtering on log operations
type LogOptions struct {
	filter.LabelMatch
	Follow bool
}
