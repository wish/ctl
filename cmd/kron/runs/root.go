package runs

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetRunsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runs",
		Short: "Subcommand on recent runs of a cron job",
		Long: `Operate on the jobs started by a cron job
	Has a bunch of subcommand just like kron`,
	}

	cmd.AddCommand(GetDescribeCmd(c))
	cmd.AddCommand(GetGetCmd(c))
	cmd.AddCommand(GetLogsCmd(c))

	return cmd
}
