package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	v1 "k8s.io/api/core/v1"
	"os"
	"os/exec"
	"strings"
)

// func cpTest() {
// 	fmt.Printf("test")
// 	return
// }
const (
	// ToPod is the command line argument for copying files into the pod
	ToPod string = "in"
	// FromPod is the command line argument for copying files out of the pod
	FromPod string = "out"
)

func cpCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cp in/out APPNAME SOURCE [flags]",
		Short: "Shortcut tool to using kubectl cp",
		Long: `Shortcut tool to using kubectl cp.
Use 'cp in' to copy a file from your local machine into the pod.
Use 'cp out' to copy a file out of the pod your local machine.
For custom pods not created by 'ctl up' use --custom-pod flag.
If there are multiple pods with the same name, it will take the first pod it finds.
If no container is set, it will use the first one.`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			container, _ := cmd.Flags().GetString("container")
			out, _ := cmd.Flags().GetString("out")
			user, _ := cmd.Flags().GetString("user")

			inOrOut := args[0]
			appName := args[1]
			source := args[2]

			// Get hostname to use in job name if not supplied
			if user == "" {
				var err error
				user, err = os.Hostname()
				if err != nil {
					return errors.New("Unable to get hostname of machine")
				}
			}

			// Replace periods with dashes and convert to lower case to follow K8's name constraints
			user = strings.Replace(user, ".", "-", -1)
			user = strings.ToLower(user)
			podName := fmt.Sprintf("%s-%s", appName, user)
			lm, _ := parsing.LabelMatch(fmt.Sprintf("name=%s", podName))
			options := client.ListOptions{LabelMatch:lm}

			pod, _, _, err := c.FindAdhocPodAndAppDetails(appName, options)
			if err != nil {
				return err
			}
			if pod == nil {
				return fmt.Errorf("No pod found, try running `ctl up %s` to start your pod", appName)
			}

			podPhase := pod.Status.Phase
			// Check to see if pod is running
			if podPhase == v1.PodPending {
				return fmt.Errorf("Pod %s is still being created", pod.Name)
			}

			// Build command to pass kubectl
			outputFiles := out
			sourceFiles := source
			if inOrOut == ToPod {
				fmt.Printf("\nCopying files into POD %s in NAMESPACE %s and CONTEXT %s\n", pod.Name, pod.Namespace, pod.Context)

				// Point destination to the pod: <namespace>/<pod>:<destination directory>
				outputFiles = fmt.Sprintf("%s/%s:%s", pod.Namespace, pod.Name, out)
			} else if inOrOut == FromPod {
				fmt.Printf("\nCopying files out of POD %s in NAMESPACE %s and CONTEXT %s\n", pod.Name, pod.Namespace, pod.Context)

				// Set source from the pod: <namespace>/<pod>:<source directory>
				sourceFiles = fmt.Sprintf("%s/%s:%s", pod.Namespace, pod.Name, source)
			} else {
				return errors.New("please use in or out to copy files/directories to and from the pods respectively")
			}
			context := fmt.Sprintf("--context=%s", pod.Context)
			if container != "" {
				container = fmt.Sprintf("--container=%s", container)
			}

			// Print out useful info for users
			fmt.Printf("\nCopying files from: %s\n", source)
			fmt.Printf("Placing them in: %v\n", out)
			command := exec.Command("kubectl", "cp", sourceFiles, outputFiles, context, container)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Stdin = os.Stdin

			// Build command
			return command.Run()
		},
	}

	cmd.Flags().StringP("out", "o", "/tmp/ctl", "Specify output folder, default to /tmp/ctl")
	cmd.Flags().StringP("container", "c", "", "Specify the container")
	cmd.Flags().Bool("custom-pod", false, "Default false. If true, will find a pod with name instead of searching for pods created by ctl up")
	cmd.Flags().StringP("user", "u", "", "Name that is used for ad hoc jobs. Defaulted to hostname.")
	cmd.Flags().StringP("status", "s", "", "Filter pods by specified status if custom-pod flag is set")

	return cmd
}
