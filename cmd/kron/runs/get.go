package runs

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

// kron/getCmd represents the kron/list command
func init() {
	RunsCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get cronjob [flags]",
	Short: "Get a list of runs of a cron job",
	Long: `Get a list of runs of a cron job.
Only operates on a single cron job.
If multiple cron jobs matches the parameters, only the first is used.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		list, err := client.GetDefaultConfigClient().
			ListRunsOfCronJob(ctxs, namespace, args[0], client.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		printRunList(list)
	},
}
