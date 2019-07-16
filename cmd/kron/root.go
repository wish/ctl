package kron

import (
	"github.com/wish/ctl/cmd/kron/runs"
	"github.com/wish/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetKronCmd(c *client.Client) *cobra.Command {
	kron := &cobra.Command{
		Use:   "kron",
		Short: "A tool for cron on kubernetes",
		Long:  "A subcommand for managing and reviewing cron jobs on kubernetes.",
	}

	kron.AddCommand(GetDescribeCmd(c))
	kron.AddCommand(GetExecCmd(c))
	kron.AddCommand(GetFavoriteCmd(c))
	kron.AddCommand(GetGetCmd(c))
	kron.AddCommand(GetSelectCmd(c))
	kron.AddCommand(GetSuspendCmd(c))
	kron.AddCommand(GetUnfavoriteCmd(c))
	kron.AddCommand(GetUnsuspendCmd(c))
	kron.AddCommand(GetWebCmd(c))
	kron.AddCommand(runs.GetRunsCmd(c))

	return kron
}
