package config

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
)

// Cmd returns the config subcommand
func Cmd(c *client.Client) *cobra.Command {
	config := &cobra.Command{
		Use:   "config",
		Short: "Edit ctl configuration",
		Long:  "Tool for changing the behaviour of ctl",
	}

	config.AddCommand(fetchCmd(c))
	config.AddCommand(deleteCmd(c))
	config.AddCommand(viewCmd(c))

	return config
}
