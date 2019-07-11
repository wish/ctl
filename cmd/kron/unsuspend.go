package kron

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetUnsuspendCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "unsuspend cronjob [flags]",
		Short: "Unsuspend a cron job",
		Long: `Unsuspends the specified cron job.
	If the cron job is not suspended, does nothing.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")

			success, err := c.SetCronJobSuspend(ctxs, namespace, args[0], false)
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
