package cmd

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringSliceP("context", "c", []string{}, "Specify the context")
	describeCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
}

var describeCmd = &cobra.Command{
	Use:   "describe pods [flags]",
	Short: "Show details of a specific pod(s)",
	Long: `Print a detailed description of the pods specified by name.
If namespace not specified, it will get all the pods across all the namespaces.
If context(s) not specified, it will search through all contexts.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		pods, err := client.GetDefaultConfigClient().FindPods(ctxs, namespace, args, client.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		if len(pods) == 0 {
			fmt.Println("Could not find any matching pods!")
		} else {
			describePodList(pods)
		}
	},
}
