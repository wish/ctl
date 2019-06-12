package cmd

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(shCmd)
	shCmd.Flags().StringSliceP("context", "c", []string{}, "Specify the context")
	shCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
	shCmd.Flags().StringP("container", "t", "", "Specify the container")
}

var shCmd = &cobra.Command{
	Use:   "sh pod [flags]",
	Short: "Exec $SHELL into the container of a specific pod",
	Long: `Exec shell into the container of a specific pod.
If namespace not specified, it will get all the pods across all the namespaces.
If context(s) not specified, it will list from all contexts.
Note that this command only operates on one pod, if multiple pods match,
the command will only work on the first one found.
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
