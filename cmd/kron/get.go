package kron

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

// kron/getCmd represents the kron/list command
func init() {
	KronCmd.AddCommand(getCmd)
	// Contexts flag
	getCmd.Flags().StringSliceP("context", "c", []string{}, "Specific contexts to list cronjobs from")
	getCmd.Flags().StringP("namespace", "n", "", "Specific namespaces to list cronjobs from within contexts")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a list of cronjobs",
	Long:  "Get a list of cronjobs based on specified search criteria.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		list, err := client.GetDefaultConfigClient().
			ListCronJobsOverContexts(ctxs, namespace, client.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		if len(list) == 0 {
			fmt.Println("No cron jobs found!")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "NAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tAGE\tCONTEXT")

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
