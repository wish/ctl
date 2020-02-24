package cron

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func unsuspendCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "unsuspend cronjob [flags]",
		Short: "Unsuspend a cron job",
		Long: `Unsuspends the specified cron job.
	If the cron job is not suspended, does nothing.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd, args)
			if err != nil {
				return err
			}

			all, err := c.ListCronJobsOverContexts(ctxs, namespace, options)
			if err != nil {
				return err
			}
			if len(all) == 0 {
				return errors.New("no cronjobs found")
			} else if len(all) > 1 {
				return errors.New("too many cronjobs match the criteria")
			}

			success, err := c.SetCronJobSuspend(all[0].Context, all[0].Namespace, all[0].Name, false)
			if err != nil {
				return err
			}

			if success {
				cmd.Println("Successfully unsuspended cron job", args[0])
			} else {
				cmd.Printf("Cron job \"%s\" was already unsuspended\n", args[0])
			}
			return err
		},
	}
}
