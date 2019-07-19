package cmd

import (
	"fmt"
	"github.com/wish/ctl/pkg/client/types"
	"os"
	"strings"
	"text/tabwriter"
	"time"
	// "gopkg.in/yaml.v2"
)

// REVIEW: Most of the processing here was guessed with reverse engineering
// by comparing with the output of kubectl
func printPodList(lst []types.PodDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No pods found")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "CONTEXT\tNAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")
	for _, l := range labelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Context)
		fmt.Fprintf(w, "\t%s", v.Namespace)
		fmt.Fprintf(w, "\t%s", v.Name)
		var ready int
		for _, s := range v.Status.ContainerStatuses {
			if s.Ready {
				ready++
			}
		}
		fmt.Fprintf(w, "\t%d/%d", ready, len(v.Spec.Containers))
		fmt.Fprintf(w, "\t%s", v.Status.Phase) // A bit off from kubectl output
		// Restarts
		var restarts int32
		for _, s := range v.Status.ContainerStatuses {
			restarts += s.RestartCount
		}
		fmt.Fprintf(w, "\t%d", restarts)
		fmt.Fprintf(w, "\t%v", time.Since(v.CreationTimestamp.Time).Round(time.Second))

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

func describePod(pod types.PodDiscovery) {
	fmt.Printf("Context: %s\n", pod.Context)
	// fmt.Printf("Namespace: %s\n", pod.Namespace)
	// b, _ := yaml.Marshal(pod.Pod)

	b, s := pod.Pod.Descriptor()

	fmt.Println(string(b))
	fmt.Println(s)
}

func describePodList(lst []types.PodDiscovery) {
	for _, pod := range lst {
		describePod(pod)
	}
}
