package runs

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func GetGetCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "get cronjob [flags]",
		Short: "Get a list of runs of a cron job",
		Long: `Get a list of runs of a cron job.
	Only operates on a single cron job.
	If multiple cron jobs matches the parameters, only the first is used.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			ctxs, err := cmd.Flags().GetStringSlice("context")
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			namespace, _ := cmd.Flags().GetString("namespace")

			list, err := c.ListRunsOfCronJob(ctxs, namespace, args[0], client.ListOptions{})

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			printRunList(list)
		},
	}
}
