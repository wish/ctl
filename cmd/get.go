package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"strings"
)

var supportedGetTypes = [][]string{
	{"pods", "pod", "po"},
	{"jobs", "job"},
	{"configmaps", "configmap", "cm"},
	{"deployments", "deployment", "deploy"},
	{"replicasets", "replicaset", "rs"},
	{"cronjobs", "cronjob"},
}

func getResourceStr() string {
	var str strings.Builder

	fmt.Fprintln(&str, "Choose from the list of supported resources:")
	for _, names := range supportedGetTypes {
		fmt.Fprintf(&str, " * %s\n", names[0])
	}

	return str.String()
}

func getCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get RESOURCE [NAME...] [flags]",
		Short: "Prints a table of resources of a type",
		Long: `Prints out a table of resource type matching the query.
Optionally, it filters through names match any of the regular expressions set.` + "\n\n" + getResourceStr(),
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

			labelColumns, _ := cmd.Flags().GetStringSlice("label-columns")

			switch strings.ToLower(args[0]) {
			case "pods", "pod", "po":
				list, err := c.ListPodsOverContexts(ctxs, namespace, options)
				// NOTE: List is unsorted and could be in an inconsistent order
				// Output
				if err != nil {
					printPodList(cmd.OutOrStdout(), list, labelColumns)
					return err
				}
				printPodList(cmd.OutOrStdout(), list, labelColumns)
			case "jobs", "job":
				list, err := c.ListJobsOverContexts(ctxs, namespace, options)
				if err != nil {
					printJobList(cmd.OutOrStdout(), list, labelColumns)
					return err
				}
				printJobList(cmd.OutOrStdout(), list, labelColumns)
			case "configmaps", "configmap", "cm":
				list, err := c.ListConfigMapsOverContexts(ctxs, namespace, options)
				if err != nil {
					printConfigMapList(cmd.OutOrStdout(), list, labelColumns)
					return err
				}
				printConfigMapList(cmd.OutOrStdout(), list, labelColumns)
			case "deployments", "deployment", "deploy":
				list, err := c.ListDeploymentsOverContexts(ctxs, namespace, options)
				if err != nil {
					printDeploymentList(cmd.OutOrStdout(), list, labelColumns)
					return err
				}
				printDeploymentList(cmd.OutOrStdout(), list, labelColumns)
			case "replicasets", "replicaset", "rs":
				list, err := c.ListReplicaSetsOverContexts(ctxs, namespace, options)
				if err != nil {
					printReplicaSetList(cmd.OutOrStdout(), list, labelColumns)
					return err
				}
				printReplicaSetList(cmd.OutOrStdout(), list, labelColumns)
			case "cronjobs", "cronjob":
				list, err := c.ListCronJobsOverContexts(ctxs, namespace, options)
				if err != nil {
					printCronJobList(cmd.OutOrStdout(), list, labelColumns)
					return err
				}
				printCronJobList(cmd.OutOrStdout(), list, labelColumns)
			case "k8s_env":
				printK8sEnvList(c.K8Envs)
			default:
				cmd.Help()
				return errors.New(`the resource type "` + args[0] + `" was not found`)
			}
			return nil
		},
	}

	cmd.Flags().StringSlice("label-columns", nil, "Prints with columns that contain the value of the specified label")

	return cmd
}
