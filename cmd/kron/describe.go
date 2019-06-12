package kron

import (
	"fmt"

	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	KronCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringSliceP("context", "c", []string{}, "Specific contexts to search cronjobs from")
	describeCmd.Flags().StringP("namespace", "n", "", "Specific namespaces to search cronjobs from within contexts")
	describeCmd.Flags().BoolP("favorites", "f", false, "Describe all favorited cron jobs")
}

// Currently does not support selected job
// Requires job name
var describeCmd = &cobra.Command{
	Use:   "describe [jobs] [flags]",
	Short: "Show details about specified cron jobs",
	Long: `Show details about specific cron jobs, or the selected job if none is specified.
If namespace not specified, it will get all the cron jobs across all the namespaces.
If context(s) not specified, it will search through all contexts.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")
		onlyFavorites, _ := cmd.Flags().GetBool("favorites")

		var cronjobs []client.CronJobDiscovery
		var err error

		if onlyFavorites {
			cronjobs, err = client.GetDefaultConfigClient().ListCronJobsOverContexts(ctxs, namespace, client.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			cronjobs = filterFromFavorites(cronjobs)
		} else if len(args) == 0 { // Use selected
			selected, err := getSelected()
			if err != nil {
				panic(err.Error())
			}
			cronjobs, err = client.GetDefaultConfigClient().FindCronJobs(selected.Location.Contexts, selected.Location.Namespace, []string{selected.Name}, client.ListOptions{})
		} else {
			cronjobs, err = client.GetDefaultConfigClient().
				FindCronJobs(ctxs, namespace, args, client.ListOptions{})
		}

		if err != nil {
			panic(err.Error())
		}

		for _, cronjob := range cronjobs {
			describeCronJob(cronjob)
		}

		if len(cronjobs) == 0 {
			fmt.Println("Did not find any matching jobs")
		}
	},
}
