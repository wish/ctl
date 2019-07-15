package cmd

import (
	"github.com/ContextLogic/ctl/cmd/util/parsing"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			shell, _ := cmd.Flags().GetString("shell")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}

			err = c.ExecInPod(ctxs, namespace, args[0], container, options, []string{shell}, os.Stdin, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")
	cmd.Flags().StringP("shell", "s", "/bin/bash", "Specify the shell path")

	return cmd
}
