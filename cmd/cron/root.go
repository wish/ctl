package cron

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
)

// Cmd returns the cron subcommand given a client to operate on
func Cmd(c *client.Client) *cobra.Command {
	cron := &cobra.Command{
		Use:   "cron",
		Short: "A tool for cron on kubernetes",
		Long:  "A subcommand for managing and reviewing cron jobs on kubernetes.",
	}

	cron.AddCommand(execCmd(c))
	cron.AddCommand(suspendCmd(c))
	cron.AddCommand(unsuspendCmd(c))

	return cron
}
