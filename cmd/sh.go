package cmd

import (
  "github.com/ContextLogic/ctl/pkg/client"
  "github.com/spf13/cobra"
  "os"
)

func init() {
  rootCmd.AddCommand(shCmd)
  shCmd.Flags().StringSliceP("namespace", "n", []string{}, "Specify the namespace")
  shCmd.Flags().StringP("context", "c", "", "Specify the context")
  shCmd.Flags().StringP("container", "t", "", "Specify the container")
}

var shCmd = &cobra.Command{
  Use: "sh pod [flags]",
  Short: "Exec $SHELL into the container of a specific pod",
  Long: `Exec shell into the container of a specific pod.
    If the pod has only one container, the container name is optional.
    If the pod has multiple containers, user have to choose one from them.`,
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    ctxs, _ := cmd.Flags().GetStringSlice("context")
    namespace, _ := cmd.Flags().GetString("namespace")
    container, _ := cmd.Flags().GetString("container")

    err := client.GetDefaultConfigClient().ExecInPod(ctxs, namespace, args[0], container, []string{"/bin/sh"}, os.Stdin, os.Stdout, os.Stderr)
    if err != nil {
      panic(err.Error())
    }
  },
}
