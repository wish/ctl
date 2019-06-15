package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	KronCmd.AddCommand(suspendCmd)
}

var suspendCmd = &cobra.Command{
	Use:   "suspend cronjob [flags]",
	Short: "Suspend a cron job",
	Long: `Suspends the specified cron job.
If the cron job is already suspended, does nothing.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		success, err := client.GetDefaultConfigClient().
			SetCronJobSuspend(ctxs, namespace, args[0], true)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if success {
			fmt.Println("Successfully suspended cron job", args[0])
		} else {
			fmt.Printf("Cron job \"%s\" was already suspended\n", args[0])
		}
	},
}
