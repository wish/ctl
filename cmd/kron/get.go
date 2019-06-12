package kron

import (
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
)

// kron/getCmd represents the kron/list command
func init() {
	KronCmd.AddCommand(getCmd)
	// Contexts flag
	getCmd.Flags().StringSliceP("context", "c", []string{}, "Specific contexts to list cronjobs from")
	getCmd.Flags().StringP("namespace", "n", "", "Specific namespaces to list cronjobs from within contexts")
	getCmd.Flags().BoolP("favorites", "f", false, "Get all favorited cron jobs")
}

var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Get a list of cronjobs",
	Long: `Get a list of cron jobs in the specified namespace and context(s).
If namespace not specified, it will get all the cron jobs across all the namespaces.
If context(s) not specified, it will list from all contexts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")
		onlyFavorites, _ := cmd.Flags().GetBool("favorites")

		list, err := client.GetDefaultConfigClient().
			ListCronJobsOverContexts(ctxs, namespace, client.ListOptions{})

		if err != nil {
			panic(err.Error())
		}

		if onlyFavorites {
			list = filterFromFavorites(list)
		}

		printCronJobList(list)
	},
}
