package cmd

import (
	"fmt"

	"os"

	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(shCmd)

	shCmd.Flags().StringP("container", "c", "", "Sepcify the container")
	shCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
	shCmd.Flags().StringP("shell", "s", "/bin/bash", "Specify the shell path")

}

var shCmd = &cobra.Command{
	Use:   "sh [pod] [flags]",
	Short: "Exec /bin/bash into the container of a specific pod",
	Long: `Print the logs for a container in a pod or specified resource. If the pod has only one container, the container name is
optional. If the pod has multiple containers, user have to choose one from them.`,
	Run: func(cmd *cobra.Command, args []string) {

		container, _ := cmd.Flags().GetString("container")
		namespace, _ := cmd.Flags().GetString("namespace")
		shell, _ := cmd.Flags().GetString("shell")

		shPod(args[0], container, namespace, shell)
	},
	Args: cobra.MinimumNArgs(1),
}

func shPod(pod, container, namespace, shell string) {
	pods, err := findPods(pod, namespace)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Printf("failed to find pod \"%s\"\n", pod)
		os.Exit(1)
	}

	if len(podSelected.Containers) > 1 {
		option := selector(podSelected.Containers)
		container = podSelected.Containers[option]
	} else if len(podSelected.Containers) == 1 {
		container = podSelected.Containers[0]
	} else {
		fmt.Printf("Pod: %s has no container.\n", podSelected.Name)
	}

	command := exec.Command("kubectl", "exec",
		"-it", podSelected.Name,
		"-c", container, shell,
		"-n", podSelected.Namespace,
		"--context", podSelected.Cluster)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if viper.GetBool("verbose") {
		prettyPrintCmd(command)
	}

	command.Run()
}
