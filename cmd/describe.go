package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"strings"
)

var supportedDescribeTypes = [][]string{
	{"pods", "pod", "po"},
	{"configmaps", "configmap", "cm"},
	{"deployments", "deployment", "deploy"},
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
		Use:   "describe pods [NAME...] [flags]",
		Short: "Show details of a specific pod(s)",
		Long: `Print a detailed description of the pods specified by name.
The names are regex expressions. ` + "\n\n" + describeResourceStr(),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")

			if len(args) == 0 {
				cmd.Help()
				return errors.New("no resource type provided")
			}

			options, err := parsing.ListOptions(cmd, args[1:])
			if err != nil {
				return err
			}

			switch args[0] {
			case "pods", "pod", "po":
				pods, err := c.ListPodsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(pods) == 0 {
					return errors.New("could not find any matching pods")
				}
				describePodList(pods)
			case "configmaps", "configmap", "cm":
				configmaps, err := c.ListConfigMapsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(configmaps) == 0 {
					return errors.New("could not find any matching configmaps")
				}
				describeConfigMapList(configmaps)
			case "deployments", "deployment", "deploy":
				deployments, err := c.ListDeploymentsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(deployments) == 0 {
					return errors.New("could not find any matching deployments")
				}
				describeDeploymentList(deployments)
			default:
				cmd.Help()
				return errors.New("the resource type \"" + args[0] + "\" was not found.\nSee 'ctl describe'")
			}
			return nil
		},
	}

	return cmd
}
