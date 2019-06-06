package kron

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ContextLogic/ctl/pkg/client"
	clienthelper "github.com/ContextLogic/ctl/pkg/client/helper"
	"github.com/spf13/cobra"
)

// kron/listCmd represents the kron/list command
func init() {
	KronCmd.AddCommand(listCmd)
	// Contexts flag
	listCmd.Flags().StringSliceP("contexts", "c", clienthelper.GetContexts(), "Specific contexts to list cronjobs from")
	// Limit flag
	listCmd.Flags().Int64P("limit", "l", 0, "Limit the number of returned cron jobs")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of cronjobs",
	Long:  "Get a list of cronjobs based on specified search criteria.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		ctxs, _ := cmd.Flags().GetStringSlice("contexts")
		// Limit
		limit, _ := cmd.Flags().GetInt64("limit")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "NAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tAGE\tCONTEXT")

		list, err := client.GetDefaultConfigClient().
			ListCronJobsOverContexts(ctxs, "", client.ListOptions{limit})

		if err != nil {
			panic(err.Error())
		}

		for _, v := range list {
			fmt.Fprintf(w, "%s\t", v.Name)          // Name
			fmt.Fprintf(w, "%s\t", v.Spec.Schedule) // Schedule
			fmt.Fprintf(w, "%t\t", *v.Spec.Suspend) // Suspend
			fmt.Fprintf(w, "%d\t", len(v.Status.Active))
			// Last schedule
			// TODO fix rounding
			if v.Status.LastScheduleTime == nil {
				fmt.Fprintf(w, "<none>\t")
			} else {
				fmt.Fprintf(w, "%v\t", time.Since(v.Status.LastScheduleTime.Time).Round(time.Second))
			}
			// Age
			fmt.Fprintf(w, "%v\t", time.Since(v.CreationTimestamp.Time).Round(time.Second))
			// Context
			fmt.Fprintf(w, "%s\n", v.Context)
		}
		w.Flush()
	},
}
