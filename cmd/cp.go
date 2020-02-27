package cmd

import (
	"fmt"
	"os/exec"
	"os"
	"errors"

	"github.com/wish/ctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/parsing"

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
		Use:   "cp in/out POD SOURCE [flags]",
		Short: "Shortcut tool to using kubectl cp",
		Long: `Shortcut tool to using kubectl cp.
Use 'cp in' to copy a file from your local macine into the pod.
Use 'cp out' to copy a file out of the pod your local machine.
If there are multiple pods with the same name, it will take the first pod it finds.
If no container is set, it will use the first one.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			container, _ := cmd.Flags().GetString("container")
			out, _ := cmd.Flags().GetString("out")
			options, err := parsing.ListOptions(cmd, nil)
			if err != nil {
				return err
			}

			inOrOut := args[0]
			nameOfPod := args[1]
			source := args[2]

			// Find pod, if there are multiple pods, pick the first one
			pod, container, err := c.FindPodWithContainer(ctxs, namespace, nameOfPod, container, options)
			if err != nil {
				return err
			}

			podNamespace := pod.Namespace
			podName := pod.Name
			podContext := pod.Context

			// Build command to pass kubectl
			outputFiles := out
			sourceFiles := source
			if inOrOut == ToPod {
				fmt.Printf("\nCopying files into POD %s in NAMESPACE %s and CONTECT %s\n", podName, podNamespace, podContext)

				// Point destination to the pod: <namespace>/<pod>:<destination directory>
				outputFiles = fmt.Sprintf("%s/%s:%s", podNamespace, podName, out)
			} else if inOrOut == FromPod {
				fmt.Printf("\nCopying files out of POD %s in NAMESPACE %s and CONTECT %s\n", podName, podNamespace, podContext)

				// Set source from the pod: <namespace>/<pod>:<source directory>
				sourceFiles = fmt.Sprintf("%s/%s:%s", podNamespace, podName, source)
			} else {
				return errors.New("please use in or out to copy files/directories to and from the pods respectively")
			}
			context := fmt.Sprintf("--context=%s", podContext)
			if container != "" {
				container = fmt.Sprintf("--container=%s", container)
			}

			// Print out useful info for users
			fmt.Printf("\nCopying files from: %s\n", source)
			fmt.Printf("Placing them in: %v\n\n", out)
			fmt.Printf("Running command...\nkubectl cp %v %v %v %v\n\n", sourceFiles, outputFiles, context, container)

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

	return cmd
}
