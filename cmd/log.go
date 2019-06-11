package cmd

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
	// "io"
)

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringSliceP("context", "c", "", "Specify the context")
	logCmd.Flags().StringP("namespace", "n", []string{}, "Specify the namespace")
	logCmd.Flags().StringP("container", "t", "", "Specify the container")
}

var logCmd = &cobra.Command{
	Use:   "log pod [flags]",
	Short: "Get log of a container in a pod",
	Long: `Print a detailed description of the selected pods..
    If namespace not specified, it will get all the pods across all the namespaces.
    If context(s) not specified, it will go through all contexts.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")
		container, _ := cmd.Flags().GetString("container")

		res, err := client.GetDefaultConfigClient().LogPod(ctxs, namespace, args[0], container, client.LogOptions{})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			// panic(err.Error())
		}

		raw, err := res.Raw()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(raw))
	},
}
