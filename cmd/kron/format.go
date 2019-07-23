package kron

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/wish/ctl/pkg/client/types"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func overrideFavoriteMessage(job string, val location) string {
	return fmt.Sprintf(`Overriding previous entry of "%s"
    Contexts: %v\n
    Namespace: %v\n`, job, val.Contexts, val.Namespace)
}

func printCronJobList(lst []types.CronJobDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No cron jobs found!")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tNEXT RUN\tAGE")
	for _, l := range labelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Name)
		fmt.Fprintf(w, "\t%s", v.Spec.Schedule) // Schedule
		fmt.Fprintf(w, "\t%t", *v.Spec.Suspend) // Suspend
		fmt.Fprintf(w, "\t%d", len(v.Status.Active))
		// Last schedule
		// TODO fix rounding
		if v.Status.LastScheduleTime == nil {
			fmt.Fprintf(w, "\t<none>")
		} else {
			fmt.Fprintf(w, "\t%v", time.Since(v.Status.LastScheduleTime.Time).Round(time.Second))
		}
		// Next run
		s, _ := cron.ParseStandard(v.Spec.Schedule)
		fmt.Fprintf(w, "\t%v", time.Until(s.Next(time.Now())).Round(time.Second))
		// Age
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))
		// Labels
		for _, l := range labelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeCronJob(c types.CronJobDiscovery) {
	fmt.Printf("context: %s\n", c.Context)
	b, _ := yaml.Marshal(c.CronJob)
	fmt.Println(string(b))
}
