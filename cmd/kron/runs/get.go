package runs

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func getCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get cronjob [flags]",
		Short: "Get a list of jobs belonging to a cron job",
		Long: `Get a list of jobs belonging to a cron job.
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
			options, err := parsing.ListOptions(cmd, nil)
			if err != nil {
				return err
			}

			list, err := c.ListJobsOfCronJob(ctxs, namespace, args[0], options)

			if err != nil {
				return err
			}

			labelColumns, _ := cmd.Flags().GetStringSlice("label-columns")
			printJobList(list, labelColumns)
			return nil
		},
	}

	cmd.Flags().StringSlice("label-columns", nil, "Prints with columns that contain the value of the specified label")

	return cmd
}
