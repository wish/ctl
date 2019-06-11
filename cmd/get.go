package cmd

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringSliceP("context", "c", []string{}, "Specify the context")
	getCmd.Flags().StringP("namespace", "n", "", "Specify the namespace")
}

var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get a list of pods",
	Long: `Get a list of pods in specified namespace and context(s).
    If namespace not specified, it will get all the pods across all the namespaces.
    If context(s) not specified, it will go through all contexts.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		list, err := client.GetDefaultConfigClient().
			ListPodsOverContexts(ctxs, namespace, client.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		// Output
		printPodList(list)
	},
}
