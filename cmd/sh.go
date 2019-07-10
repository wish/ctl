package cmd

import (
	"fmt"
	"os"

	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetShCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
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
			ctxs, err := cmd.Flags().GetStringSlice("context")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			shell, _ := cmd.Flags().GetString("shell")

			err = c.ExecInPod(ctxs, namespace, args[0], container, []string{shell}, os.Stdin, os.Stdout, os.Stderr)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")
	cmd.Flags().StringP("shell", "s", "/bin/bash", "Specify the shell path")

	return cmd
}
