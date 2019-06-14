package cmd

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringP("container", "t", "", "Specify the container")
}

var logsCmd = &cobra.Command{
	Use:     "logs pod [flags]",
	Aliases: []string{"log"},
	Short:   "Get log of a container in a pod",
	Long: `Print a detailed description of the selected pod.
If namespace not specified, it will get all the pods across all the namespaces.
If context(s) not specified, it will search through all contexts.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")
		container, _ := cmd.Flags().GetString("container")

		res, err := client.GetDefaultConfigClient().LogPodOverContexts(ctxs, namespace, args[0], container, client.LogOptions{})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			// panic(err.Error())
		}

		raw, err := res.Raw()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// REVIEW: Format??
		fmt.Println(string(raw))
	},
}
