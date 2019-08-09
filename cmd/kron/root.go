package kron

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
)

// Cmd returns the kron subcommand given a client to operate on
func Cmd(c *client.Client) *cobra.Command {
	kron := &cobra.Command{
		Use:   "kron",
		Short: "A tool for cron on kubernetes",
		Long:  "A subcommand for managing and reviewing cron jobs on kubernetes.",
	}

	kron.AddCommand(execCmd(c))
	kron.AddCommand(suspendCmd(c))
	kron.AddCommand(unsuspendCmd(c))
	kron.AddCommand(webCmd(c))

	return kron
}
