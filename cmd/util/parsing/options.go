package parsing

import (
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/spf13/cobra"
)

func LabelMatchFromCmd(cmd *cobra.Command) (filter.LabelMatch, error) {
	s, _ := cmd.Flags().GetStringArray("label")
	return LabelMatchSlice(s)
}

func ListOptions(cmd *cobra.Command) (client.ListOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.ListOptions{l}, err
}

func GetOptions(cmd *cobra.Command) (client.GetOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.GetOptions{l}, err
}

func LogOptions(cmd *cobra.Command) (client.LogOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.LogOptions{l}, err
}
