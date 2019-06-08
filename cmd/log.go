package cmd

import (
  "github.com/ContextLogic/ctl/pkg/client"
  "github.com/spf13/cobra"
  "fmt"
  // "os"
  // "io"
)

func init() {
  rootCmd.AddCommand(logCmd)
  logCmd.Flags().StringSliceP("namespace", "n", []string{}, "Specify the namespace")
  logCmd.Flags().StringP("context", "c", "", "Specify the context")
  logCmd.Flags().StringP("container", "t", "", "Specify the container")
}

var logCmd = &cobra.Command{
  Use: "log pod [flags]",
  Short: "Get log of a container in a pod",
  Long: `Print a detailed description of the selected pods..
    If namespace not specified, it will get all the pods across all the namespaces.
    If context(s) not specified, it will go through all contexts.`,
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    ctxs, _ := cmd.Flags().GetStringSlice("context")
    namespace, _ := cmd.Flags().GetString("namespace")
    container, _ := cmd.Flags().GetString("container")

    err := client.GetDefaultConfigClient().LogPod(ctxs, namespace, args[0], container, client.LogOptions{})
    if err != nil {
      fmt.Println(err.Error())
      // panic(err.Error())
    }
    // _, err = io.Copy(os.Stdout, readCloser)
    // readCloser.Close()
    // pods, err := client.GetDefaultConfigClient().FindPods(ctxs, namespaces, args, client.ListOptions{})
    // if err != nil {
    //   panic(err.Error())
    // }
    // if len(pods) == 0 {
    //   fmt.Println("Could not find any matching pods!")
    // } else {
    //   describePodList(pods)
    // }
  },
}
