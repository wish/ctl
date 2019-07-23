package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/types"
	"strings"
)

var supportedDescribeTypes = [][]string{
	{"pods", "pod", "po"},
}

func describeResourceStr() string {
	var str strings.Builder

	fmt.Fprintln(&str, "Choose from the list of supported resources:")
	for _, names := range supportedGetTypes {
		fmt.Fprintf(&str, " * %s\n", names[0])
	}

	return str.String()
}

func describeCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe pods [flags]",
		Short: "Show details of a specific pod(s)",
		Long:  "Print a detailed description of the pods specified by name.\n\n" + describeResourceStr(),
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
				var pods []types.PodDiscovery
				var err error
				if len(args) == 1 {
					pods, err = c.ListPodsOverContexts(ctxs, namespace, options)
				} else {
					pods, err = c.FindPods(ctxs, namespace, args[1:], options)
				}
				if err != nil {
					return err
				}
				if len(pods) == 0 {
					return errors.New("could not find any matching pods")
				}
				describePodList(pods)
			default:
				defer cmd.Help()
				return errors.New(`The resource type "` + args[0] + `" was not found.
See 'ctl describe'`)
			}
			return nil
		},
	}

	return cmd
}
