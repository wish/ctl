package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolP("follow", "f", false, "stream pod logs (stdout)")
	logCmd.Flags().StringP("container", "c", "", "Print the logs of this container")
	logCmd.Flags().StringP("tail", "t", "", "Print the logs of this container")
	logCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")

}

var logCmd = &cobra.Command{
	Use:   "log [pod] [flags]",
	Short: "Get log of a container in a pod",
	Long: `Print the logs for a container in a pod or specified resource. If the pod has only one container, the container name is
optional. If the pod has multiple containers, user have to choose one from them.`,
	Run: func(cmd *cobra.Command, args []string) {
		// compile flags
		var flags []string
		if cmd.Flags().Changed("follow") {
			flags = append(flags, "-f")
		}
		if cmd.Flags().Changed("tail") {
			lines, _ := cmd.Flags().GetString("tail")
			lineNum, err := strconv.ParseInt(lines, 10, 64)
			if err != nil {
				fmt.Println("Please sepecify a valid line number.")
				os.Exit(1)
			}
			if lineNum < 0 {
				fmt.Println("tailLines must be greate than or equal to 0")
				os.Exit(1)
			}
			flags = append(flags, "--tail", lines)
		}

		container, _ := cmd.Flags().GetString("container")
		namespace, _ := cmd.Flags().GetString("namespace")

		err := logPod(args[0], container, namespace, flags)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"logs"},
}

func logPod(pod, container, namespace string, flags []string) error {
	pods, err := findPods(pod, namespace)
	if err != nil {
		return err
	}
	var podSelected Pod
	if len(pods) > 1 {
		var options []string
		for _, p := range pods {
			options = append(options, fmt.Sprintf("%s/%s/%s", p.Cluster, p.Namespace, p.Name))
		}
		podSelected = pods[selector(options)]
	} else if len(pods) == 1 {
		podSelected = pods[0]
	} else {
		return errors.Errorf("failed to find pod \"%s\"\n", pod)
	}

	if container == "" && len(podSelected.Containers) > 1 {
		option := selector(podSelected.Containers)
		container = podSelected.Containers[option]
	}

	args := []string{"logs", podSelected.Name,
		"-n", podSelected.Namespace,
		"--context", podSelected.Cluster,
		"-c", container}
	args = append(args, flags...)

	command := exec.Command("kubectl", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if viper.GetBool("verbose") {
		prettyPrintCmd(command)
	}

	return command.Run()
}
