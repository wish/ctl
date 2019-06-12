package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/robfig/cron"
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
	fmt.Fprintln(w, "NAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tNEXT RUN\tAGE\tCONTEXT")

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
		// Next run
		s, _ := cron.ParseStandard(v.Spec.Schedule)
		fmt.Fprintf(w, "%v\t", time.Until(s.Next(time.Now())).Round(time.Second))
		// Age
		fmt.Fprintf(w, "%v\t", time.Since(v.CreationTimestamp.Time).Round(time.Second))
		// Context
		fmt.Fprintf(w, "%s\n", v.Context)
	}
	w.Flush()
}

func describeCronJob(c client.CronJobDiscovery) {
	fmt.Printf("Context: %s\n", c.Context)
	fmt.Printf("\tName: %s\n", c.Name)
	fmt.Printf("\tNamespace: %s\n", c.Namespace)
	fmt.Printf("\tSchedule: %s\n", c.Spec.Schedule)
	fmt.Printf("\tActive: %d\n", len(c.Status.Active))
	fmt.Printf("\tLast schedule: %v\n", time.Since(c.Status.LastScheduleTime.Time).Round(time.Second))
	s, _ := cron.ParseStandard(c.Spec.Schedule)
	fmt.Printf("\tNext run: %v\n", time.Until(s.Next(time.Now())).Round(time.Second))
	fmt.Printf("\tCreated on: %v\n", c.CreationTimestamp)
}
