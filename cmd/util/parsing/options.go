package parsing

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/filter"
)

// LabelMatchFromCmd automatically parses the "label" flag from a command
// and returns the filtering.LabelMatch specified.
func LabelMatchFromCmd(cmd *cobra.Command) (filter.LabelMatch, error) {
	s, _ := cmd.Flags().GetStringArray("label")
	return LabelMatchSlice(s)
}

// ListOptions parses a client.ListOptions from a command
func ListOptions(cmd *cobra.Command) (client.ListOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.ListOptions{l}, err
}

// GetOptions parses a client.GetOptions from a command
func GetOptions(cmd *cobra.Command) (client.GetOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.GetOptions{l}, err
}

// LogOptions parses a client.LogOptions from a command
func LogOptions(cmd *cobra.Command) (client.LogOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.LogOptions{l}, err
}
