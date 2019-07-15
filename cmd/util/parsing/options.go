package parsing

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/ContextLogic/ctl/pkg/client/filter"
	"github.com/spf13/cobra"
)

func labelMatchFromCmd(cmd *cobra.Command) (filter.LabelMatch, error) {
	s, _ := cmd.Flags().GetStringArray("labels")
	return LabelMatchSlice(s)
}

func ListOptions(cmd *cobra.Command) (client.ListOptions, error) {
	l, err := labelMatchFromCmd(cmd)
	return client.ListOptions{l}, err
}

func GetOptions(cmd *cobra.Command) (client.GetOptions, error) {
	l, err := labelMatchFromCmd(cmd)
	return client.GetOptions{l}, err
}

func LogOptions(cmd *cobra.Command) (client.LogOptions, error) {
	l, err := labelMatchFromCmd(cmd)
	return client.LogOptions{l}, err
}
