package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

var supportedGetTypes = [][]string{
	{"pods", "pod", "po"},
}

func getCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get RESOURCE [flags]",
		Short: "Get a list of resources",
		Long: `Get a list of reources in the specified namespace and context(s).
If namespace not specified, it will get all the resources across all the namespaces.
If context(s) not specified, it will list from all contexts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				defer cmd.Help()
				return errors.New("no resource type provided")
			}

			switch args[0] {
			case "pods", "pod", "po":
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
				defer cmd.Help()
				return errors.New(`The resource type "` + args[0] + `" was not found`)
			}
			return nil
		},
	}

	cmd.Flags().StringSlice("label-columns", nil, "Prints with columns that contain the value of the specified label")

	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		cmd.Println("Choose from the list of supported resources:")
		for _, names := range supportedGetTypes {
			cmd.Printf(" * %s\n", names[0])
		}
		return nil
	})

	return cmd
}
