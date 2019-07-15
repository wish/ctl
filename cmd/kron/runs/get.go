package runs

import (
	"github.com/ContextLogic/ctl/cmd/util/parsing"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetGetCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "get cronjob [flags]",
		Short: "Get a list of runs of a cron job",
		Long: `Get a list of runs of a cron job.
Only operates on a single cron job.
If multiple cron jobs matches the parameters, only the first is used.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			ctxs, err := cmd.Flags().GetStringSlice("context")
			if err != nil {
				return err
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}

			list, err := c.ListRunsOfCronJob(ctxs, namespace, args[0], options)

			if err != nil {
				return err
			}

			printRunList(list)
			return nil
		},
	}
}
