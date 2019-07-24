package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func logsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs pod [flags]",
		Aliases: []string{"log"},
		Short:   "Get log of a container in a pod",
		Long:    `Print a detailed description of the selected pod.`,
		Args:    cobra.ExactArgs(1),
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
			cmd.Println(string(raw))

			return nil
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")

	return cmd
}
