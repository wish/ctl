package runs

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
)

// Cmd returns the kron/runs subcommand given a client
func Cmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runs",
		Short: "Subcommand on recent runs of a cron job",
		Long: `Operate on the jobs started by a cron job
Has a bunch of subcommand just like kron`,
	}

	cmd.AddCommand(describeCmd(c))
	cmd.AddCommand(getCmd(c))
	cmd.AddCommand(logsCmd(c))

	return cmd
}
