package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	// "io"
	"bufio"
)

func logsCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs pod [flags]",
		Aliases: []string{"log"},
		Short:   "Get log of a container in a pod",
		Long:    `Print a detailed description of the selected pod.`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			follow, _ := cmd.Flags().GetBool("follow")
			options, err := parsing.LogOptions(cmd, args)
			options.Follow = follow
			if err != nil {
				return err
			}

			if follow {
				reader, err := c.LogPodsOverContexts(ctxs, namespace, container, options)
				if err != nil {
					return err
				}
				scanner := bufio.NewScanner(reader)

				for scanner.Scan() {
					cmd.Println(scanner.Text())
				}

				if err = scanner.Err(); err != nil {
					return err
				}
			} else {
				req, err := c.LogPodOverContexts(ctxs, namespace, args[0], container, options)
				res := req.Do()
				raw, err := res.Raw()
				if err != nil {
					return err
				}
				cmd.Print(string(raw))
			}

			return nil
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")
	cmd.Flags().BoolP("follow", "f", false, "Specify if the logs should be streamed")

	return cmd
}
