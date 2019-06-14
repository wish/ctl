package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	KronCmd.AddCommand(unsuspendCmd)
}

var unsuspendCmd = &cobra.Command{
	Use:   "unsuspend cronjob [flags]",
	Short: "Unsuspend a cron job",
	Long: `Unsuspends the specified cron job.
If the cron job is not suspended, does nothing.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		success, err := client.GetDefaultConfigClient().
			SetCronJobSuspend(ctxs, namespace, args[0], false)
		if err != nil {
			panic(err.Error())
		}

		if success {
			fmt.Println("Successfully unsuspended cron job", args[0])
		} else {
			fmt.Printf("Cron job \"%s\" was already unsuspended\n", args[0])
		}
	},
}
