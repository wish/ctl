package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"os"
)

func shCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sh POD [flags]",
		Short: "Exec $SHELL into the container of a specific pod",
		Long: `Exec shell into the container of a specific pod.
Note that this command only operates on one pod, if multiple pods have the exact name,
the command will only work on the first one found.
If the pod has only one container, the container name is optional.
If the pod has multiple containers, user have to choose one from them.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			shell, _ := cmd.Flags().GetString("shell")
			options, err := parsing.ListOptions(cmd, nil)
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
