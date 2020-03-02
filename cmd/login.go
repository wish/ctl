package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	v1 "k8s.io/api/core/v1"
)

func loginCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login APPNAME [flags]",
		Short: "Exec $SHELL into the pod of your ad hoc job spawned by ctl up",
		Long: `Exec shell into the container of the pod spawned by ctl up.
Note that this command only operates on one pod, if multiple pods have the exact name,
the command will only work on the first one found.
If the pod has only one container, the container name is optional.
If the pod has multiple containers, it will choose the first container found.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			shell, _ := cmd.Flags().GetString("shell")
			user, _ := cmd.Flags().GetString("user")

			appName := args[0]

			// Get hostname to use in job name if not supplied
			if user == "" {
				var err error
				user, err = os.Hostname()
				if err != nil {
					return errors.New("Unable to get hostname of machine")
				}
			}

			// We get the pod through the name label
			podName := fmt.Sprintf("%s-%s", appName, user)
			lm, _ := parsing.LabelMatch(fmt.Sprintf("name=%s", podName))
			options := client.ListOptions{LabelMatch: lm}

			pods, err := c.ListPodsOverContexts(ctxs, namespace, options)
			if err != nil {
				return err
			}
			if len(pods) < 1 {
				return fmt.Errorf("No pod found, try running `ctl up %s` to start your pod", appName)
			}

			pod := pods[0]

			podPhase := pod.Status.Phase
			// Check to see if pod is running
			if podPhase == v1.PodPending {
				return fmt.Errorf("Pod %s is still being created", pod.Name)
			}

			fmt.Printf("Shelling in pod: %s...\n"+
				"Use `ctl cp in %s <files>` to copy files into pod\n"+
				"Use `ctl cp out %s <files>` to copy files out of pod\n"+
				"Use `ctl cp -h` for more info about file copying\n\n",
				pod.Name, pod.Name, pod.Name)

			err = c.ExecInPod([]string{pod.Context}, pod.Namespace, pod.Name, container, client.ListOptions{},
				[]string{shell}, os.Stdin, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringP("container", "c", "", "Specify the container")
	cmd.Flags().StringP("shell", "s", "/bin/bash", "Specify the shell path")
	cmd.Flags().StringP("user", "u", "", "Name that is used for ad hoc jobs. Defaulted to hostname.")

	return cmd
}
