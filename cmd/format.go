package cmd

import (
  "fmt"
  "github.com/ContextLogic/ctl/pkg/client"
  "os"
  "text/tabwriter"
  "time"
)

func printPodList(lst []client.PodDiscovery) {
  w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
  fmt.Fprintln(w, "CONTEXT\tNAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")

  for _, v := range lst {
    fmt.Fprintf(w, "%s\t", v.Context)
    fmt.Fprintf(w, "%s\t", v.Namespace)
    fmt.Fprintf(w, "%s\t", v.Name)
    fmt.Fprintf(w, "TODO\t")
    fmt.Fprintf(w, "%v\t", v.Status.Phase)
    fmt.Fprintf(w, "TODO\t")
    fmt.Fprintf(w, "%v\n", time.Since(v.CreationTimestamp.Time).Round(time.Second))
  }
  w.Flush()
}

func describePod(pod client.PodDiscovery) {
  fmt.Printf("Context: %s\n", pod.Context)
  fmt.Printf("Namespace: %s\n", pod.Namespace)
  fmt.Printf("%+v\n", pod.Pod)
}

func describePodList(lst []client.PodDiscovery) {
  for _, pod := range lst {
    describePod(pod)
  }
}
