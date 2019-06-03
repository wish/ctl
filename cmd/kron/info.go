package kron

import (
	"fmt"
	"time"
	"github.com/ContextLogic/wishctl/pkg/kron"
	"github.com/spf13/cobra"
)

func init() {
	KronCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringSliceP("contexts", "c", kron.GetContexts(), "Specific contexts to list cronjobs from")
	infoCmd.Flags().StringSliceP("namespaces", "n", []string{}, "Specific namespaces to list cronjobs from within contexts")
}

// Currently does not support selected job
// Requires job name
var infoCmd = &cobra.Command{
	Use:   "info [job]",
	Short: "Get info about a job",
	Long:  "Get info about a specific job, or the selected job if none is specified.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Attempting to find job \"%s\"\n", args[0])

		// Contexts
		ctxs, _ := cmd.Flags().GetStringSlice("contexts")
		// Namespaces
		nss, _ := cmd.Flags().GetStringSlice("namespaces")
		// Positional arg
		job := args[0]

		for _, ctx := range ctxs {
			cl, err := kron.GetContextClient(ctx)
			if err != nil {
				fmt.Printf("ERROR: Context \"%s\" not found\n", ctx)
				continue;
			}

			var namespaces []string
			if len(nss) == 0 {
				namespaces = cl.GetNamespaces()
			} else {
				namespaces = nss
			}

			for _, ns := range namespaces {
				cronjob, err := cl.Get(ns, job, kron.GetOptions{})
				if err != nil {
					// Cronjob not found on this context
					continue;
				}

				fmt.Printf("Context: %s\n", ctx)
				fmt.Printf("\tNamespace: %s\n", ns)
				fmt.Printf("\tSchedule: %s\n", cronjob.Spec.Schedule)
				fmt.Printf("\tActive: %d\n", len(cronjob.Status.Active))
				fmt.Printf("\tLast schedule: %v\n", time.Since(cronjob.Status.LastScheduleTime.Time).Round(time.Second))
				fmt.Printf("\tCreated on: %v\n", cronjob.CreationTimestamp)
			}
		}
	},
}
