package cmd

import (
	"fmt"
	"os"

	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get a list of pods",
	Long: `Get a list of pods in the specified namespace and context(s).
If namespace not specified, it will get all the pods across all the namespaces.
If context(s) not specified, it will list from all contexts.`,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := client.GetDefaultConfigClient().
			ListPodsOverContexts(ctxs, namespace, client.ListOptions{})
		// NOTE: List is unsorted and could be in an inconsistent order
		// Output
		if list != nil {
			printPodList(list)
		}
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}
