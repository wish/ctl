package cmd

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client/types"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

// REVIEW: Most of the processing here was guessed with reverse engineering
// by comparing with the output of kubectl
func printPodList(w io.Writer, lst []types.PodDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Fprintln(w, "No pods found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(tw, "CONTEXT\tNAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(tw, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(tw)

	for _, v := range lst {
		fmt.Fprintf(tw, "%s", v.Context)
		fmt.Fprintf(tw, "\t%s", v.Namespace)
		fmt.Fprintf(tw, "\t%s", v.Name)
		var ready int
		for _, s := range v.Status.ContainerStatuses {
			if s.Ready {
				ready++
			}
		}
		fmt.Fprintf(tw, "\t%d/%d", ready, len(v.Spec.Containers))
		fmt.Fprintf(tw, "\t%s", v.Status.Phase) // A bit off from kubectl output
		// Restarts
		var restarts int32
		for _, s := range v.Status.ContainerStatuses {
			restarts += s.RestartCount
		}
		fmt.Fprintf(tw, "\t%d", restarts)
		fmt.Fprintf(tw, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(tw, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(tw, v.Labels[l])
			}
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

func printConfigMapList(w io.Writer, lst []types.ConfigMapDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Fprintln(w, "No config maps found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(tw, "CONTEXT\tNAMESPACE\tNAME\tDATA\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(tw, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(tw)

	for _, v := range lst {
		fmt.Fprintf(tw, "%s", v.Context)
		fmt.Fprintf(tw, "\t%s", v.Namespace)
		fmt.Fprintf(tw, "\t%s", v.Name)
		fmt.Fprintf(tw, "\t%d", len(v.Data))
		// Age
		fmt.Fprintf(tw, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(tw, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(tw, v.Labels[l])
			}
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

func printDeploymentList(w io.Writer, lst []types.DeploymentDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Fprintln(w, "No deployments found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(tw, "CONTEXT\tNAMESPACE\tNAME\tDESIRED\tCURRENT\tUP-TO-DATE\tAVAILABLE\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(tw, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(tw)

	for _, v := range lst {
		fmt.Fprintf(tw, "%s", v.Context)
		fmt.Fprintf(tw, "\t%s", v.Namespace)
		fmt.Fprintf(tw, "\t%s", v.Name)
		// Desired
		fmt.Fprintf(tw, "\t%d", *v.Spec.Replicas)
		// Current
		fmt.Fprintf(tw, "\t%d", v.Status.Replicas)
		// Up-to-date
		fmt.Fprintf(tw, "\t%d", v.Status.UpdatedReplicas)
		// Available
		fmt.Fprintf(tw, "\t%d", v.Status.AvailableReplicas)
		// Age
		fmt.Fprintf(tw, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(tw, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(tw, v.Labels[l])
			}
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

func printReplicaSetList(w io.Writer, lst []types.ReplicaSetDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Fprintln(w, "No replicasets found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(tw, "CONTEXT\tNAMESPACE\tNAME\tDESIRED\tCURRENT\tREADY\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(tw, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(tw)

	for _, v := range lst {
		fmt.Fprintf(tw, "%s", v.Context)
		fmt.Fprintf(tw, "\t%s", v.Namespace)
		fmt.Fprintf(tw, "\t%s", v.Name)
		// Desired
		fmt.Fprintf(tw, "\t%d", *v.Spec.Replicas)
		// Current
		fmt.Fprintf(tw, "\t%d", v.Status.Replicas)
		// Ready
		fmt.Fprintf(tw, "\t%d", v.Status.ReadyReplicas)
		// Age
		fmt.Fprintf(tw, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

		for _, l := range newLabelColumns {
			fmt.Fprint(tw, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(tw, v.Labels[l])
			}
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

func printJobList(w io.Writer, lst []types.JobDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Fprintln(w, "No replicasets found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(tw, "CONTEXT\tNAMESPACE\tNAME\tSTATE\tSTART\tEND")
	for _, l := range newLabelColumns {
		fmt.Fprint(tw, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(tw)

	for _, v := range lst {
		fmt.Fprintf(tw, "%s", v.Context)
		fmt.Fprintf(tw, "\t%s", v.Namespace)
		fmt.Fprintf(tw, "\t%s", v.Name)
		// State
		if v.Status.Failed > 0 {
			fmt.Fprint(tw, "\tFAILED")
		} else if v.Status.CompletionTime != nil {
			fmt.Fprint(tw, "\tSUCCESSFUL")
		} else {
			fmt.Fprint(tw, "\tIN PROGRESS")
		}
		// Start
		fmt.Fprintf(tw, "\t%v", v.Status.StartTime)
		// END
		if v.Status.CompletionTime != nil {
			fmt.Fprintf(tw, "\t%v", v.Status.CompletionTime)
		} else {
			fmt.Fprint(tw, "\t<none>")
		}

		for _, l := range newLabelColumns {
			fmt.Fprint(tw, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(tw, v.Labels[l])
			}
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

func printCronJobList(w io.Writer, lst []types.CronJobDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Fprintln(w, "No cronjobs found")
		return
	}
	// Insert default columns
	defaultColumns := viper.GetStringSlice("default_columns")
	var newLabelColumns []string
	if len(defaultColumns) == 0 {
		newLabelColumns = labelColumns
	} else if len(labelColumns) == 0 {
		newLabelColumns = defaultColumns
	} else {
		for _, s := range defaultColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
		for _, s := range labelColumns {
			newLabelColumns = append(newLabelColumns, s)
		}
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(tw, "CONTEXT\tNAME\tSCHEDULE\tSUSPEND\tACTIVE\tLAST SCHEDULE\tNEXT RUN\tAGE")
	for _, l := range newLabelColumns {
		fmt.Fprint(tw, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(tw)

	for _, v := range lst {
		fmt.Fprintf(tw, "%s", v.Context)
		fmt.Fprintf(tw, "\t%s", v.Name)
		fmt.Fprintf(tw, "\t%s", v.Spec.Schedule) // Schedule
		fmt.Fprintf(tw, "\t%t", *v.Spec.Suspend) // Suspend
		fmt.Fprintf(tw, "\t%d", len(v.Status.Active))
		// Last schedule
		// TODO fix rounding
		if v.Status.LastScheduleTime == nil {
			fmt.Fprintf(tw, "\t<none>")
		} else {
			fmt.Fprintf(tw, "\t%v", time.Since(v.Status.LastScheduleTime.Time).Round(time.Second))
		}
		// Next run
		s, _ := cron.ParseStandard(v.Spec.Schedule)
		fmt.Fprintf(tw, "\t%v", time.Until(s.Next(time.Now())).Round(time.Second))
		// Age
		fmt.Fprintf(tw, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))
		// Labels
		for _, l := range newLabelColumns {
			fmt.Fprint(tw, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(tw, v.Labels[l])
			}
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
}

func printK8sEnvList(k8sEnv []string) {
	for _, v := range k8sEnv {
		fmt.Println(v)
	}
}
