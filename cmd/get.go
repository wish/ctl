package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
	getCmd.Flags().StringP("region", "r", "", "Specify the region")
	getCmd.Flags().StringP("env", "e", "", "Specify the enviroment")
	getCmd.Flags().StringP("az", "a", "", "Specify the alvalibility zone")
	getCmd.Flags().StringP("config", "", "", "Specify the config file")

}

var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get a list of pods",
	Long:  "Get a list of pods in namespace. If namespace not specified, it will get all the pods across all the namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		region, _ := cmd.Flags().GetString("region")
		env, _ := cmd.Flags().GetString("env")
		az, _ := cmd.Flags().GetString("az")
		config, _ := cmd.Flags().GetString("config")
		getPodsFromClustersByNamespace(namespace, config, region, env, az)
	},
	Args:    cobra.MaximumNArgs(0),
	Aliases: []string{"gets"},
}

func getPodsInNamespaceInCluster(cluster, namespace string, resultChannel chan PodList) {

	var tempPodList PodList
	var command *exec.Cmd

	if namespace == "" {
		command = exec.Command(
			"kubectl",
			"get", "pods",
			"--all-namespaces",
			//"-o=custom-columns=NAME:.metadata.name",
			fmt.Sprintf("--context=%s", cluster))
	} else {
		command = exec.Command(
			"kubectl",
			"get", "pods",
			"-n", namespace,
			//"-o=custom-columns=NAME:.metadata.name",
			fmt.Sprintf("--context=%s", cluster))
	}

	if viper.GetBool("verbose") {
		prettyPrintCmd(command)
	}

	result, err := command.Output()

	if err != nil || string(result) == "" {
		resultChannel <- tempPodList
		return
	}

	podList := strings.Split(strings.TrimSpace(string(result)), "\n")
	if len(podList) < 2 {
		resultChannel <- tempPodList
		return
	}

	resultChannel <- PodList{cluster, podList}
}

func getPodsFromClustersByNamespace(namespace, configpath, region, enviroment, az string) {
	clusters, err := getFilteredClusters(configpath, region, enviroment, az)
	if err != nil {
		fmt.Printf("failed to get clusters: %v\n", err)
		os.Exit(1)
	}
	if len(clusters) == 0 {
		fmt.Print("failed to get clusters: no cluster found\n")
		os.Exit(1)
	}

	resultChan := make(chan PodList)

	for _, c := range clusters {
		go getPodsInNamespaceInCluster(c, namespace, resultChan)
	}

	var result []PodList

	for range clusters {
		response := <-resultChan
		if len(response.Pods) == 0 {
			continue
		}

		fmt.Print(response)
		result = append(result, response)
	}

	if len(result) == 0 {
		fmt.Println("No pod found.")
	}
}
