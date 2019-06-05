package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().BoolP("follow", "f", false, "stream pod logs (stdout)")
	logCmd.Flags().StringP("container", "c", "", "Print the logs of this container")
	logCmd.Flags().StringP("tail", "t", "", "lines of most recent log to be printed")
	logCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
	logCmd.Flags().BoolP("aggregate", "", false, "Aggregated all the logs found ")
	logCmd.Flags().StringP("since", "s", "", "Only return logs newer than a relative duration like 5s, 2m, or 3h")
	logCmd.Flags().StringP("region", "r", "", "Specify the region")
	logCmd.Flags().StringP("env", "e", "", "Specify the enviroment")
	logCmd.Flags().StringP("az", "a", "", "Specify the alvalibility zone")
	logCmd.Flags().StringP("config", "", "", "Specify the config file")
	logCmd.Flags().String("since-time", "", "Only return logs after a specific date (RFC3339). Defaults to all logs")

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
		if cmd.Flags().Changed("since") {
			duration, _ := cmd.Flags().GetString("since")
			flags = append(flags, "--since", duration)
		}
		if cmd.Flags().Changed("since-time") {
			sinceTime, _ := cmd.Flags().GetString("since-time")
			flags = append(flags, "--since-time", sinceTime)
		}

		container, _ := cmd.Flags().GetString("container")
		namespace, _ := cmd.Flags().GetString("namespace")
		aggregate, _ := cmd.Flags().GetBool("aggregate")
		region, _ := cmd.Flags().GetString("region")
		env, _ := cmd.Flags().GetString("env")
		az, _ := cmd.Flags().GetString("az")
		config, _ := cmd.Flags().GetString("config")

		logPod(args[0], container, namespace, config, region, env, az, aggregate, flags)
	},
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"logs"},
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("follow") && cmd.Flags().Changed("aggregate") {
			fmt.Println("Cannot aggregate logs while streaming logs")
			os.Exit(1)
		}
		if cmd.Flags().Changed("since") && cmd.Flags().Changed("since-time") {
			fmt.Println("Only one of since-time / since may be used")
			os.Exit(1)
		}
	},
}

func logPod(pod, container, namespace, configpath, region, environment, az string, aggregate bool, flags []string) {
	pods, err := findPods(pod, namespace, configpath, region, environment, az)
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

	for _, p := range podSelected {
		logSinglePod(p, container, flags)
	}

}

func logSinglePod(pod Pod, container string, flags []string) {
	if container == "" {
		container = pod.Containers[0]
	}

	fmt.Printf("Printing log from %s/%s/%s - %s:\n", pod.Cluster, pod.Namespace, pod.Name, container)
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
