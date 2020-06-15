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
	{"jobs", "job"},
	{"configmaps", "configmap", "cm"},
	{"deployments", "deployment", "deploy"},
	{"replicasets", "replicaset", "rs"},
	{"cronjobs", "cronjob"},
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
			showEvents, _ := cmd.Flags().GetBool("show-events")

			if len(args) == 0 {
				cmd.Help()
				return errors.New("no resource type provided")
			}

			options, err := parsing.ListOptions(cmd, args[1:])
			if err != nil {
				return err
			}

			describeOptions := client.DescribeOptions{ShowEvents: showEvents}

			switch strings.ToLower(args[0]) {
			case "pods", "pod", "po":
				pods, err := c.ListPodsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(pods) == 0 {
					return errors.New("could not find any matching pods")
				}
				for _, pod := range pods {
					s, err := c.DescribePod(pod, describeOptions)
					if err != nil {
						cmd.Println(err.Error())
					} else {
						cmd.Println(s)
					}
				}
			case "jobs", "job":
				jobs, err := c.ListJobsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(jobs) == 0 {
					return errors.New("could not find any matching jobs")
				}
				for _, job := range jobs {
					s, err := c.DescribeJob(job, describeOptions)
					if err != nil {
						cmd.Println(err.Error())
					} else {
						cmd.Println(s)
					}
				}
			case "configmaps", "configmap", "cm":
				configmaps, err := c.ListConfigMapsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(configmaps) == 0 {
					return errors.New("could not find any matching configmaps")
				}
				for _, configmap := range configmaps {
					s, err := c.DescribeConfigMap(configmap, describeOptions)
					if err != nil {
						cmd.Println(err.Error())
					} else {
						cmd.Println(s)
					}
				}
			case "deployments", "deployment", "deploy":
				deployments, err := c.ListDeploymentsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(deployments) == 0 {
					return errors.New("could not find any matching deployments")
				}
				for _, deployment := range deployments {
					s, err := c.DescribeDeployment(deployment, describeOptions)
					if err != nil {
						cmd.Println(err.Error())
					} else {
						cmd.Println(s)
					}
				}
			case "replicasets", "replicaset", "rs":
				replicasets, err := c.ListReplicaSetsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(replicasets) == 0 {
					return errors.New("could not find any matching replicasets")
				}
				for _, replicaset := range replicasets {
					s, err := c.DescribeReplicaSet(replicaset, describeOptions)
					if err != nil {
						cmd.Println(err.Error())
					} else {
						cmd.Println(s)
					}
				}
			case "cronjobs", "cronjob":
				cronjobs, err := c.ListCronJobsOverContexts(ctxs, namespace, options)
				if err != nil {
					return err
				}
				if len(cronjobs) == 0 {
					return errors.New("could not find any matching cronjobs")
				}
				for _, cronjob := range cronjobs {
					s, err := c.DescribeCronJob(cronjob, describeOptions)
					if err != nil {
						cmd.Println(err.Error())
					} else {
						cmd.Println(s)
					}
				}
			default:
				cmd.Help()
				return errors.New("the resource type \"" + args[0] + "\" was not found.\nSee 'ctl describe'")
			}
			return nil
		},
	}

	cmd.Flags().Bool("show-events", true, "If true, display events related to the described object.")
	cmd.Flags().StringP("status", "s", "", "Filter pods by specified status")

	return cmd
}
