package kron

import (
	"fmt"
	"time"
	"sync"
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

		var waitc sync.WaitGroup
		waitc.Add(len(ctxs))

		for _, ctx := range ctxs {
			go func(ctx string) {
				defer waitc.Done()

				cl, err := kron.GetContextClient(ctx)
				if err != nil {
					fmt.Printf("ERROR: Context \"%s\" not found\n", ctx)
					return
				}

				var namespaces []string
				if len(nss) == 0 {
					namespaces = cl.GetNamespaces()
				} else {
					namespaces = nss
				}

				var waitn sync.WaitGroup
				waitn.Add(len(namespaces))

				for _, ns := range namespaces {
					go func(ns string) {
						defer waitn.Done()

						cronjob, err := cl.Get(ns, job, kron.GetOptions{})
						if err != nil {
							// Cronjob not found on this context
							return
						}

						fmt.Printf("Context: %s\n\tNamespace: %s\n\tSchedule: %s\n\tActive: %d\n\tLast Schedule: %v\n\tCreated on: %v\n",
							ctx, ns, cronjob.Spec.Schedule, len(cronjob.Status.Active), time.Since(cronjob.Status.LastScheduleTime.Time).Round(time.Second), cronjob.CreationTimestamp)
					} (ns)
				}
				waitn.Wait()
			} (ctx)
		}
		waitc.Wait()
	},
}
