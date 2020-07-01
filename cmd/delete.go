package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
)

var supportedDeleteTypes = [][]string{
	{"pods", "pod", "po"},
	{"jobs", "job"},
	{"configmaps", "configmap", "cm"},
	{"deployments", "deployment", "deploy"},
	{"replicasets", "replicaset", "rs"},
	{"cronjobs", "cronjob"},
}

func deleteResourceStr() string {
	var str strings.Builder

	fmt.Fprintln(&str, "Choose from the list of supported resources:")
	for _, names := range supportedGetTypes {
		fmt.Fprintf(&str, " * %s\n", names[0])
	}

	return str.String()
}

// Default Y/n
func prompter(r io.Reader) bool {
	rr := bufio.NewReader(r)
	s, err := rr.ReadString('\n')
	if err != nil {
		panic(err)
	}
	s = strings.ToLower(s)[:len(s)-1]
	return s == "y" || s == "yes" || s == ""
}

func deleteCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete RESOURCE NAME... [flags]",
		Short: "Deletes resources",
		Long:  `Delete all matching resources of a type. If more than one resource is found, it prints a table and asks for confirmation.` + "\n\n" + getResourceStr(),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")

			if len(args) == 0 {
				cmd.Help()
				return errors.New("no resource type provided")
			} else if len(args) == 1 {
				cmd.Help()
				return errors.New("must specify a name for delete")
			}

			options, err := parsing.ListOptions(cmd, args[1:])
			if err != nil {
				return err
			}

			labelColumns, _ := cmd.Flags().GetStringSlice("label-columns")

			now, _ := cmd.Flags().GetBool("now")
			deleteChildren, _ := cmd.Flags().GetBool("delete-children")
			deleteOptions := client.DeleteOptions{Now: now, DeletionPropagation: deleteChildren}

			switch strings.ToLower(args[0]) {
			case "pods", "pod", "po":
				list, err := c.ListPodsOverContexts(ctxs, namespace, options)
				// Output
				if err != nil {
					return err
				}

				if len(list) > 1 { // Warn about deleting
					printPodList(cmd.OutOrStdout(), list, labelColumns)
					cmd.Printf("\nAre you sure you want to delete the %d items above? [Y/n] ", len(list))
					if !prompter(cmd.InOrStdin()) {
						cmd.Println("Aborted")
						return nil
					}
				}

				errstr := ""
				count := 0
				for _, l := range list {
					err := c.DeletePod(l.Context, l.Namespace, l.Name, deleteOptions)
					if err != nil {
						errstr += err.Error()
						count++
					}
				}
				if count > 0 {
					cmd.Printf("Encountered %d errors while deleting pods", count)
					return errors.New(errstr)
				}
				cmd.Printf("Deleted %d pods\n", len(list))

			case "jobs", "job":
				list, err := c.ListJobsOverContexts(ctxs, namespace, options)
				// Output
				if err != nil {
					return err
				}

				if len(list) > 1 { // Warn about deleting
					printJobList(cmd.OutOrStdout(), list, labelColumns)
					cmd.Printf("\nAre you sure you want to delete the %d items above? [Y/n] ", len(list))
					if !prompter(cmd.InOrStdin()) {
						cmd.Println("Aborted")
						return nil
					}
				}

				errstr := ""
				count := 0
				for _, l := range list {
					err := c.DeleteJob(l.Context, l.Namespace, l.Name, deleteOptions)
					if err != nil {
						errstr += err.Error()
						count++
					}
				}
				if count > 0 {
					cmd.Printf("Encountered %d errors while deleting jobs", count)
					return errors.New(errstr)
				}
				cmd.Printf("Deleted %d jobs\n", len(list))

			case "configmaps", "configmap", "cm":
				list, err := c.ListConfigMapsOverContexts(ctxs, namespace, options)
				// Output
				if err != nil {
					return err
				}

				if len(list) > 1 { // Warn about deleting
					printConfigMapList(cmd.OutOrStdout(), list, labelColumns)
					cmd.Printf("\nAre you sure you want to delete the %d items above? [Y/n] ", len(list))
					if !prompter(cmd.InOrStdin()) {
						cmd.Println("Aborted")
						return nil
					}
				}

				errstr := ""
				count := 0
				for _, l := range list {
					err := c.DeleteConfigMap(l.Context, l.Namespace, l.Name, deleteOptions)
					if err != nil {
						errstr += err.Error()
						count++
					}
				}
				if count > 0 {
					cmd.Printf("Encountered %d errors while deleting configmaps", count)
					return errors.New(errstr)
				}
				cmd.Printf("Deleted %d configmaps\n", len(list))

			case "deployments", "deployment", "deploy":
				list, err := c.ListDeploymentsOverContexts(ctxs, namespace, options)
				// Output
				if err != nil {
					return err
				}

				if len(list) > 1 { // Warn about deleting
					printDeploymentList(cmd.OutOrStdout(), list, labelColumns)
					cmd.Printf("\nAre you sure you want to delete the %d items above? [Y/n] ", len(list))
					if !prompter(cmd.InOrStdin()) {
						cmd.Println("Aborted")
						return nil
					}
				}

				errstr := ""
				count := 0
				for _, l := range list {
					err := c.DeleteDeployment(l.Context, l.Namespace, l.Name, deleteOptions)
					if err != nil {
						errstr += err.Error()
						count++
					}
				}
				if count > 0 {
					cmd.Printf("Encountered %d errors while deleting deployments", count)
					return errors.New(errstr)
				}
				cmd.Printf("Deleted %d deployments\n", len(list))

			case "replicasets", "replicaset", "rs":
				list, err := c.ListReplicaSetsOverContexts(ctxs, namespace, options)
				// Output
				if err != nil {
					return err
				}

				if len(list) > 1 { // Warn about deleting
					printReplicaSetList(cmd.OutOrStdout(), list, labelColumns)
					cmd.Printf("\nAre you sure you want to delete the %d items above? [Y/n] ", len(list))
					if !prompter(cmd.InOrStdin()) {
						cmd.Println("Aborted")
						return nil
					}
				}

				errstr := ""
				count := 0
				for _, l := range list {
					err := c.DeleteReplicaSet(l.Context, l.Namespace, l.Name, deleteOptions)
					if err != nil {
						errstr += err.Error()
						count++
					}
				}
				if count > 0 {
					cmd.Printf("Encountered %d errors while deleting replicasets", count)
					return errors.New(errstr)
				}
				cmd.Printf("Deleted %d replicasets\n", len(list))
			case "cronjobs", "cronjob":
				list, err := c.ListCronJobsOverContexts(ctxs, namespace, options)
				// Output
				if err != nil {
					return err
				}

				if len(list) > 1 { // Warn about deleting
					printCronJobList(cmd.OutOrStdout(), list, labelColumns)
					cmd.Printf("\nAre you sure you want to delete the %d items above? [Y/n] ", len(list))
					if !prompter(cmd.InOrStdin()) {
						cmd.Println("Aborted")
						return nil
					}
				}

				errstr := ""
				count := 0
				for _, l := range list {
					err := c.DeleteCronJob(l.Context, l.Namespace, l.Name, deleteOptions)
					if err != nil {
						errstr += err.Error()
						count++
					}
				}
				if count > 0 {
					cmd.Printf("Encountered %d errors while deleting cronjobs", count)
					return errors.New(errstr)
				}
				cmd.Printf("Deleted %d cronjobs\n", len(list))

			default:
				cmd.Help()
				return errors.New(`the resource type "` + args[0] + `" was not found`)
			}
			return nil
		},
	}

	cmd.Flags().StringSlice("label-columns", nil, "Prints with columns that contain the value of the specified label")
	cmd.Flags().Bool("now", false, "If true, signals the resource for immediate shutdown")
	cmd.Flags().Bool("delete-children", true, "If true, deletes jobs' spawned resources")
	cmd.Flags().StringP("status", "s", "", "Filter pods by specified status")

	return cmd
}
