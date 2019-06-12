package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"os"
	"text/tabwriter"
	"time"
)

func overrideFavoriteMessage(job string, val location) string {
	return fmt.Sprintf(`Overriding previous entry of "%s"
    Contexts: %v\n
    Namespace: %v\n`, job, val.Contexts, val.Namespace)
}

func printCronJobList(lst []client.CronJobDiscovery) {
	if len(lst) == 0 {
		fmt.Println("No cron jobs found!")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tAGE\tCONTEXT")

	for _, v := range lst {
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
}
