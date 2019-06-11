package kron

import (
	"fmt"
	"time"

	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	KronCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringSliceP("context", "c", []string{}, "Specific contexts to list cronjobs from")
	describeCmd.Flags().StringP("namespace", "n", "", "Specific namespaces to list cronjobs from within contexts")
}

// Currently does not support selected job
// Requires job name
var describeCmd = &cobra.Command{
	Use:   "describe [jobs]",
	Short: "Show details about specified cron jobs",
	Long:  "Show details about specific cron jobs, or the selected job if none is specified.",
	Run: func(cmd *cobra.Command, args []string) {
		// Contexts
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		// Namespaces
		namespace, _ := cmd.Flags().GetString("namespace")

		cronjobs, err := client.GetDefaultConfigClient().
			FindCronJobs(ctxs, namespace, args, client.ListOptions{})

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
