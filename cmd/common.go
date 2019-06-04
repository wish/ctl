package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// Pod defines the pod struct
type Pod struct {
	Name       string
	Namespace  string
	Cluster    string
	Containers []string
}

// WishCtlError defines errors of wishctl
type WishCtlError string

// Error implements error interfaces
func (w WishCtlError) Error() string { return string(w) }

// CLUSTERNOTFOUND defines the error of not finding any cluster
const CLUSTERNOTFOUND = WishCtlError("failed to get cluster, no cluster found")

func findPods(pod, namespace string) ([]Pod, error) {
	clusters, err := getAllClusters()
	if err != nil {
		return nil, err
	}
	if len(clusters) == 0 {
		return nil, CLUSTERNOTFOUND
	}

	resultChan := make(chan []Pod)

	for _, c := range clusters {
		go findPodsInCluster(pod, c, namespace, resultChan)
	}

	var res []Pod
	for range clusters {
		response := <-resultChan
		if response == nil {
			continue
		}
		res = append(res, response...)
	}

	return res, nil
}

func findPodsInCluster(pod, cluster, namespace string, resultChan chan []Pod) {

	var command *exec.Cmd

	if namespace == "" {
		command = exec.Command("kubectl", "get", "pods",
			"--all-namespaces", "--no-headers", "--context", cluster,
			"-o=custom-columns=NAME:.metadata.name,NAMESPACE:.metadata.namespace,CONTAINERS:.spec.containers[*].name")
	} else {
		command = exec.Command("kubectl", "get", "pods",
			"-n", namespace, "--no-headers", "--context", cluster,
			"-o=custom-columns=NAME:.metadata.name,NAMESPACE:.metadata.namespace,CONTAINERS:.spec.containers[*].name")
	}

	if viper.GetBool("verbose") {
		prettyPrintCmd(command)
	}

	result, err := command.Output()
	if err != nil || string(result) == "" {
		resultChan <- nil
		return
	}

	podList := strings.Split(strings.TrimSpace(string(result)), "\n")
	if len(podList) < 1 {
		resultChan <- nil
		return
	}

	var tmpPodList []Pod
	for _, p := range podList {
		pList := strings.Fields(p)
		name := pList[0]
		namespace := pList[1]
		containers := strings.Split(pList[2], ",")
		if strings.Contains(name, pod) {
			tmpPodList = append(tmpPodList, Pod{name, namespace, cluster, containers})
		}
	}
	resultChan <- tmpPodList
}

// selector prompts a user interface for choosing from multiple options
func selector(options []string) int {
	reader := bufio.NewReader(os.Stdin)
	for i, o := range options {
		fmt.Printf("%d:\t%s\n", i+1, o)
	}
	fmt.Printf("Select a pod/container (1-%d): ", len(options))
	input, _ := reader.ReadString('\n')
	num, err := strconv.Atoi(strings.TrimSpace(input))
	for err != nil || num < 1 || num > len(options) {
		fmt.Printf("Please enter a valid number (1-%d): ", len(options))
		input, _ = reader.ReadString('\n')
		num, err = strconv.Atoi(strings.TrimSpace(input))
	}
	return num - 1
}

func prettyPrintCmd(command *exec.Cmd) {
	var fmtCmd = "Running: "
	for _, a := range command.Args {
		fmtCmd += a + " "
	}
	color.Green(fmtCmd)
}
