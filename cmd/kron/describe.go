package kron

import (
	"fmt"
	"time"

	"github.com/ContextLogic/ctl/pkg/client"
	clienthelper "github.com/ContextLogic/ctl/pkg/client/helper"
	"github.com/spf13/cobra"
)

func init() {
	KronCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringSliceP("contexts", "c", clienthelper.GetContexts(), "Specific contexts to list cronjobs from")
	describeCmd.Flags().StringSliceP("namespaces", "n", []string{}, "Specific namespaces to list cronjobs from within contexts")
}

// Currently does not support selected job
// Requires job name
var describeCmd = &cobra.Command{
	Use:   "describe [jobs]",
	Short: "Show details about specified cron jobs",
	Long:  "Show details about specific cron jobs, or the selected job if none is specified.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Attempting to find job \"%s\"\n", args[0])

		// Contexts
		ctxs, _ := cmd.Flags().GetStringSlice("contexts")
		// Namespaces
		nss, _ := cmd.Flags().GetStringSlice("namespaces")
		// Positional arg
		job := args[0]

		cronjobs, err := client.GetDefaultConfigClient().
			GetCronJobOverMultiple(ctxs, nss, job, client.GetOptions{})

		if err != nil {
			panic(err.Error())
		}

		for _, cronjob := range cronjobs {
			// Formatter here
			fmt.Printf("Context: %s\n\tNamespace: %s\n\tSchedule: %s\n\tActive: %d\n\tLast Schedule: %v\n\tCreated on: %v\n",
				cronjob.Context, cronjob.Namespace, cronjob.Spec.Schedule, len(cronjob.Status.Active), time.Since(cronjob.Status.LastScheduleTime.Time).Round(time.Second), cronjob.CreationTimestamp)
		}

		if len(cronjobs) == 0 {
			fmt.Println("Did not find any matching jobs")
		}
	},
}
