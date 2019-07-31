package runs

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func describeCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "describe job",
		Short: "Get info about a job",
		Long:  "Get information about a specific job belonging to a cron job.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, err := cmd.Flags().GetStringSlice("command")
			if err != nil {
				return err
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd, nil)
			if err != nil {
				return err
			}

			list, err := c.FindJobs(ctxs, namespace, args, options)
			if err != nil {
				return err
			}

			if len(list) == 0 {
				return err
			}
			for _, r := range list {
				describeJob(r)
			}

			return nil
		},
	}
}
