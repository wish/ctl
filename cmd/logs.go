package cmd

import (
	"fmt"
	"github.com/ContextLogic/ctl/cmd/util/parsing"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetLogsCmd(c *client.Client) *cobra.Command {
	ret := &cobra.Command{
		Use:     "logs pod [flags]",
		Aliases: []string{"log"},
		Short:   "Get log of a container in a pod",
		Long: `Print a detailed description of the selected pod.
If namespace not specified, it will get all the pods across all the namespaces.
If context(s) not specified, it will search through all contexts.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			options, err := parsing.LogOptions(cmd)
			if err != nil {
				return err
			}

			res, err := c.LogPodOverContexts(ctxs, namespace, args[0], container, options)
			if err != nil {
				return err
			}

			raw, err := res.Raw()
			if err != nil {
				return err
			}
			// REVIEW: Format??
			fmt.Println(string(raw))

			return nil
		},
	}

	ret.Flags().StringP("container", "c", "", "Specify the container")

	return ret
}
