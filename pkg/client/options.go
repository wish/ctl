package client

import (
	"github.com/wish/ctl/pkg/client/filter"
	"regexp"
)

// ListOptions is used to specify filtering on list operations
type ListOptions struct {
	// Filtering by labels
	filter.LabelMatch
	// Filtering by name
	Search *regexp.Regexp
}

// GetOptions is used to specify filtering on get operations
type GetOptions struct {
	// Filtering by labels
	filter.LabelMatch
}

// LogOptions is used to specify filtering on log operations
type LogOptions struct {
	// Filtering by labels
	filter.LabelMatch
	// When set streams logs
	Follow bool
	// Filtering by name
	Search *regexp.Regexp
	// When set, adds a RFC3339Nano timestamp to the beginning of each line
	Timestamps bool
}

// DescribeOptions is used to specify how to describe a resource
type DescribeOptions struct {
	ShowEvents bool
}
