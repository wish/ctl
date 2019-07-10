package cmd

import (
	"fmt"
	"os"

	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func GetDescribeCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "describe pods [flags]",
		Short: "Show details of a specific pod(s)",
		Long: `Print a detailed description of the pods specified by name.
	If namespace not specified, it will get all the pods across all the namespaces.
	If context(s) not specified, it will search through all contexts.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctxs, err := cmd.Flags().GetStringSlice("context")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			namespace, _ := cmd.Flags().GetString("namespace")

			pods, err := c.FindPods(ctxs, namespace, args, client.ListOptions{})
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if len(pods) == 0 {
				fmt.Println("Could not find any matching pods!")
			} else {
				describePodList(pods)
			}
		},
	}
}
