package runs

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/ContextLogic/ctl/pkg/util"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RunsCmd.AddCommand(describeCmd)
}

var describeCmd = &cobra.Command{
	Use:   "describe run",
	Short: "Get info about a run",
	Long:  "Get information about a specific run of a cron job.", // TODO
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, err := util.GetContexts(cmd)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		namespace, _ := cmd.Flags().GetString("namespace")

		list, err := client.GetDefaultConfigClient().
			FindRuns(ctxs, namespace, args, client.ListOptions{})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if len(list) == 0 {
			fmt.Println("Could not find run")
			os.Exit(1)
		}
		for _, r := range list {
			describeRun(r)
		}
	},
}
