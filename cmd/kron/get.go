package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
	"sort"
)

func GetGetCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
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
			// Ordering of list
			l, _ := cmd.Flags().GetBool("by-last-run")
			L, _ := cmd.Flags().GetBool("by-last-run-reverse")
			e, _ := cmd.Flags().GetBool("by-next-run")
			E, _ := cmd.Flags().GetBool("by-next-run-reverse")

			if l && L || l && e || l && E || L && e || L && E || e && E { // More than one
				fmt.Println("Only at most one ordering flag may be set!")
				os.Exit(1)
			}

			list, err := c.ListCronJobsOverContexts(ctxs, namespace, client.ListOptions{})

			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			if onlyFavorites {
				list = filterFromFavorites(list)
			}

			if l {
				sort.Sort(byLastRun(list))
			} else if L {
				sort.Sort(sort.Reverse(byLastRun(list)))
			}

			printCronJobList(list)
		},
	}

	// Contexts flag
	cmd.Flags().BoolP("favorites", "f", false, "Get all favorited cron jobs")
	// Ordering flags
	cmd.Flags().BoolP("by-last-run", "l", false, "Sort chronologically by last run")
	cmd.Flags().BoolP("by-last-run-reverse", "L", false, "Sort reverse chronologically by last run")
	cmd.Flags().BoolP("by-next-run", "e", false, "Sort cronologically by next scheduled run")
	cmd.Flags().BoolP("by-next-run-reverse", "E", false, "Sort reverse chronologically by next scheduled run")

	return cmd
}
