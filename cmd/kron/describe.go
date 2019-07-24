package kron

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/types"
)

// Currently does not support selected job
// Requires job name
func describeCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe [NAME...] [flags]",
		Short: "Show details about specified cron jobs",
		Long:  `Show details about specific cron jobs, or the selected job if none is specified.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			onlyFavorites, _ := cmd.Flags().GetBool("favorites")
			options, err := parsing.ListOptions(cmd, nil)
			if err != nil {
				return err
			}

			var cronjobs []types.CronJobDiscovery

			if onlyFavorites {
				cronjobs, err = c.ListCronJobsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				cronjobs = filterFromFavorites(cronjobs)
			} else if len(args) == 0 { // Use selected
				selected, err := getSelected()
				if err != nil {
					return err
				}
				cronjobs, err = c.FindCronJobs(selected.Location.Contexts, selected.Location.Namespace, []string{selected.Name}, options)
			} else {
				cronjobs, err = c.FindCronJobs(ctxs, namespace, args, options)
			}

			if err != nil {
				return err
			}

			for _, cronjob := range cronjobs {
				describeCronJob(cronjob)
			}

			if len(cronjobs) == 0 {
				cmd.Println("Did not find any matching jobs")
			}

			return nil
		},
	}

	cmd.Flags().BoolP("favorites", "f", false, "Describe all favorited cron jobs")

	return cmd
}
