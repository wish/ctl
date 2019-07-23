package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

func getCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get RESOURCE [flags]",
		Short: "Get a list of resources",
		Long: `Get a list of reources in the specified namespace and context(s).
If namespace not specified, it will get all the resources across all the namespaces.
If context(s) not specified, it will list from all contexts.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}

			switch args[0] {
			case "pods":
				list, err := c.ListPodsOverContexts(ctxs, namespace, options)
				// NOTE: List is unsorted and could be in an inconsistent order
				// Output
				if list != nil {
					labelColumns, _ := cmd.Flags().GetStringSlice("label-columns")
					printPodList(list, labelColumns)
				}
				if err != nil {
					return err
				}
			default:
				return errors.New(`The resource type "` + args[0] + `" was not found`)
			}
			return nil
		},
	}

	cmd.Flags().StringSlice("label-columns", nil, "Prints with columns that contain the value of the specified label")

	return cmd
}
