package cmd

import (
	"errors"
	"github.com/ContextLogic/ctl/cmd/util/parsing"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetDescribeCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "describe pods [flags]",
		Short: "Show details of a specific pod(s)",
		Long: `Print a detailed description of the pods specified by name.
	If namespace not specified, it will get all the pods across all the namespaces.
	If context(s) not specified, it will search through all contexts.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}

			pods, err := c.FindPods(ctxs, namespace, args, options)
			if err != nil {
				return err
			}
			if len(pods) == 0 {
				return errors.New("Could not find any matching pods!")
			} else {
				describePodList(pods)
			}
			return nil
		},
	}
}
