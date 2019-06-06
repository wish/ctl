package cmd

import (
  "fmt"
  "github.com/ContextLogic/ctl/pkg/client"
  "text/tabwriter"
  "github.com/spf13/cobra"
  "os"
  "time"
)

func init() {
  rootCmd.AddCommand(getCmd)
  getCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
  getCmd.Flags().StringSliceP("context", "c", []string{}, "Specify the context")
}

var getCmd = &cobra.Command{
  Use: "get [flags]",
  Short: "Get a list of pods",
  Long: `Get a list of pods in specified namespace and context(s).
    If namespace not specified, it will get all the pods across all the namespaces.
    If context(s) not specified, it will go through all contexts.`,
  Run: func(cmd *cobra.Command, args []string) {
    ctxs, _ := cmd.Flags().GetStringSlice("contexts")
    namespace, _ := cmd.Flags().GetString("namespace")

    list, err := client.GetDefaultConfigClient().
      ListPodsOverContexts(ctxs, namespace, client.ListOptions{0})

    if err != nil {
      panic(err.Error())
    }

    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
    fmt.Fprintln(w, "CONTEXT\tNAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")

    for _, v := range list {
      fmt.Fprintf(w, "%s\t", v.Context)
      fmt.Fprintf(w, "%s\t", v.Namespace)
      fmt.Fprintf(w, "%s\t", v.Name)
      fmt.Fprintf(w, "TODO\t")
      fmt.Fprintf(w, "%v\t", v.Status.Phase)
      fmt.Fprintf(w, "TODO\t")
      fmt.Fprintf(w, "%v\n", time.Since(v.CreationTimestamp.Time).Round(time.Second))
    }
    w.Flush()
  },
}
