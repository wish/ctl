package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
}

var describeCmd = &cobra.Command{
	Use:   "describe [pod] [flags]",
	Short: "Show details of a specific pods",
	Long:  `Print a detailed description of the selected pods.`,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		if err := describePod(args[0], namespace); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Args: cobra.MinimumNArgs(1),
}

func describePod(pod, namespace string) error {
	pods, err := findPods(pod, namespace)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		return errors.Errorf("failed to find pod \"%s\"\n", pod)
	}

	for _, p := range pods {
		command := exec.Command("kubectl", "describe", "pods", p.Name,
			"-n", p.Namespace, "--context", p.Cluster)

		if viper.GetBool("verbose") {
			prettyPrintCmd(command)
		}

		res, err := command.Output()
		if err != nil {
			return err
		}
		fmt.Println(string(res))
	}
	return nil
}
