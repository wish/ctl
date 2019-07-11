package runs

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetDescribeCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "describe run",
		Short: "Get info about a run",
		Long:  "Get information about a specific run of a cron job.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, err := cmd.Flags().GetStringSlice("command")
			if err != nil {
				return err
			}
			namespace, _ := cmd.Flags().GetString("namespace")

			list, err := c.FindRuns(ctxs, namespace, args, client.ListOptions{})
			if err != nil {
				return err
			}

			if len(list) == 0 {
				return err
			}
			for _, r := range list {
				describeRun(r)
			}

			return nil
		},
	}
}
