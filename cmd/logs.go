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

			options, err := parsing.LogOptions(cmd, args)
			// TODO: move these options to parsing (add to kron)
			options.Follow, _ = cmd.Flags().GetBool("follow")
			options.Timestamps, _ = cmd.Flags().GetBool("timestamps")

			if err != nil {
				return err
			}

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
			return nil
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")
	cmd.Flags().BoolP("follow", "f", false, "Specify if the logs should be streamed")
	cmd.Flags().Bool("timestamps", false, "Add a RFC3339Nano format timestamp to the beginning of each line")

	return cmd
}
