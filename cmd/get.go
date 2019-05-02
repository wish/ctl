package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bytes"

	"github.com/spf13/cobra"
)

// PodList defines the list of Pods
type PodList struct {
	ClusterName string
	Pods        []string
}

// String format the PodList struct
func (p PodList) String() string {

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Cluster - %s:\n", p.ClusterName))

	if p.Pods == nil {
		b.WriteString("\tnil\n")
		return b.String()
	}

	for _, n := range p.Pods {
		b.WriteString(fmt.Sprintf("\t%s\n", n))
	}

	return b.String()

}

var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get a list of pods",
	Long:  "Get a list of pods in namespace. If namespace not specified, it will get all the pods across all the namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		getPodsFromAllClustersByNamespace(namespace)
	},
	Args:    cobra.MaximumNArgs(0),
	Aliases: []string{"gets"},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
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

func getPodsFromAllClustersByNamespace(namespace string) {
	clusters, err := getAllClusters()
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

func getAllClusters() ([]string, error) {
	result, err := exec.Command("kubectl", "config", "get-contexts", "-o=name").Output()
	if err != nil {
		return nil, err
	}
	clusterList := strings.Split(strings.TrimSpace(string(result)), "\n")
	return clusterList, nil
}
