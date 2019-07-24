package config

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/pkg/client"
)

func fetchCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Update extensions",
		Run: func(cmd *cobra.Command, args []string) {
			config.WriteCtlExt(c.GetCtlExt())
		},
	}
}
