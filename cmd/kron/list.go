package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/kron"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

// kron/listCmd represents the kron/list command

func init() {
	KronCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of cronjobs",
	Long:  "Get a list of cronjobs based on specified search criteria.",
	Run: func(cmd *cobra.Command, args []string) {
		ctxs := kron.GetContexts()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(w, "NAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tAGE\tCONTEXT")

		for _, ctx := range ctxs {
			cl, err := kron.GetContextClient(ctx)
			if err != nil {
				panic(err.Error())
			}
			list, err := cl.List()
			if err != nil {
				panic(err.Error())
			}
			for _, v := range list {
				fmt.Fprintf(w, "%s\t%s\t%t\t%d\t%v\t%v\t%s\n",
					v.Name,          // Name
					v.Spec.Schedule, // Schedule
					// Suspend boolean
					*v.Spec.Suspend,
					// Active jobs
					len(v.Status.Active),
					// Last schedule
					// TODO fix rounding
					time.Since(v.Status.LastScheduleTime.Time).Round(time.Second),
					// Age
					time.Since(v.CreationTimestamp.Time).Round(time.Second),
					// Context
					ctx)
			}
		}

		// FOR DEBUGGING the values stored in a cronjob object
		// fmt.Println("Object Meta")
		// fmt.Println(v.ObjectMeta.String())
		// fmt.Println("Spec")
		// fmt.Println(v.Spec.String())
		// fmt.Println("Status")
		// fmt.Println(v.Status.String())
		w.Flush()
	},
}
