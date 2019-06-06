package cmd

import (
  "github.com/ContextLogic/ctl/pkg/client"
  "github.com/spf13/cobra"
  "fmt"
)

func init() {
  rootCmd.AddCommand(describeCmd)
  describeCmd.Flags().StringSliceP("namespace", "n", []string{}, "Specify the namespace")
  describeCmd.Flags().StringSliceP("context", "c", []string{}, "Specify the context")
}

var describeCmd = &cobra.Command{
  Use: "describe pods [flags]",
  Short: "Show details of a specific pod(s)",
  Long: `Print a detailed description of the selected pods..
    If namespace not specified, it will get all the pods across all the namespaces.
    If context(s) not specified, it will go through all contexts.`,
  Args: cobra.MinimumNArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    ctxs, _ := cmd.Flags().GetStringSlice("context")
    namespaces, _ := cmd.Flags().GetStringSlice("namespace")

    pods, err := client.GetDefaultConfigClient().FindPods(ctxs, namespaces, args, client.ListOptions{})
    if err != nil {
      panic(err.Error())
    }
    if len(pods) == 0 {
      fmt.Println("Could not find any matching pods!")
    } else {
      describePodList(pods)
    }
  },
}
