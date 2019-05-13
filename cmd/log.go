package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolP("follow", "f", false, "stream pod logs (stdout)")
	logCmd.Flags().StringP("container", "c", "", "Print the logs of this container")
	logCmd.Flags().StringP("tail", "t", "", "lines of most recent log to be printed")
	logCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
	logCmd.Flags().BoolP("aggregate", "a", false, "Aggregated all the logs found ")

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
		aggregate, _ := cmd.Flags().GetBool("aggregate")

		logPod(args[0], container, namespace, aggregate, flags)
	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"logs"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("follow") && cmd.Flags().Changed("aggregate") {
			fmt.Println("Cannot aggregate logs while streaming logs")
			os.Exit(1)
		}
	},
}

func logPod(pod, container, namespace string, aggregate bool, flags []string) {
	pods, err := findPods(pod, namespace)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var podSelected []Pod
	if len(pods) > 1 {
		if !aggregate {
			var options []string
			for _, p := range pods {
				options = append(options, fmt.Sprintf("%s/%s/%s", p.Cluster, p.Namespace, p.Name))
			}
			podSelected = []Pod{pods[selector(options)]}
		} else {
			podSelected = pods
		}
	} else if len(pods) == 1 {
		podSelected = pods
	} else {
		fmt.Printf("failed to find pod \"%s\"\n", pod)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(podSelected))

	for _, p := range podSelected {
		go logSinglePod(p, container, flags, &wg)
	}
	wg.Wait()

}

func logSinglePod(pod Pod, container string, flags []string, wg *sync.WaitGroup) {
	defer wg.Done()
	if container == "" {
		container = pod.Containers[0]
	}

	fmt.Printf("Printing log from: %s/%s/%s - %s\n:", pod.Cluster, pod.Namespace, pod.Name, container)
	args := []string{"logs", pod.Name,
		"-n", pod.Namespace,
		"--context", pod.Cluster,
		"-c", container}
	args = append(args, flags...)

	command := exec.Command("kubectl", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if viper.GetBool("verbose") {
		prettyPrintCmd(command)
	}

	command.Run()
}
