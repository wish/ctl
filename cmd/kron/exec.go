package kron

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func execCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "exec cronjob [flags]",
		Short: "Executes a job now",
		Long: `Executes the specified job or the selected job if none.
Namespace and context flags can be set to help find the right cron job.
If multiple cron job are found, only the first one will be executed.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}

			job, err := c.RunCronJob(ctxs, namespace, args[0], options)

			if err != nil {
				return err
			}
			cmd.Printf("Successfully started job \"%s\" in context \"%s\" and namespace \"%s\"\n", job.Name, job.Context, job.Namespace)

			return nil
		},
	}
}
