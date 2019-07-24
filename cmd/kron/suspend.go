package kron

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func suspendCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "suspend cronjob [flags]",
		Short: "Suspend a cron job",
		Long: `Suspends the specified cron job.
If the cron job is already suspended, does nothing.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd, nil)
			if err != nil {
				return err
			}

			success, err := c.SetCronJobSuspend(ctxs, namespace, args[0], true, options)
			if err != nil {
				return err
			}

			if success {
				cmd.Println("Successfully suspended cron job", args[0])
			} else {
				cmd.Printf("Cron job \"%s\" was already suspended\n", args[0])
			}

			return nil
		},
	}
}
