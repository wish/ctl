package kron

import (
	"errors"
	"github.com/ContextLogic/ctl/cmd/util/parsing"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"sort"
)

func GetGetCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [flags]",
		Short: "Get a list of cronjobs",
		Long: `Get a list of cron jobs in the specified namespace and context(s).
	If namespace not specified, it will get all the cron jobs across all the namespaces.
	If context(s) not specified, it will list from all contexts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")
			onlyFavorites, _ := cmd.Flags().GetBool("favorites")
			options, err := parsing.ListOptions(cmd)
			if err != nil {
				return err
			}
			// Ordering of list
			l, _ := cmd.Flags().GetBool("by-last-run")
			L, _ := cmd.Flags().GetBool("by-last-run-reverse")
			e, _ := cmd.Flags().GetBool("by-next-run")
			E, _ := cmd.Flags().GetBool("by-next-run-reverse")

			if l && L || l && e || l && E || L && e || L && E || e && E { // More than one
				return errors.New("Only at most one ordering flag may be set!")
			}

			list, err := c.ListCronJobsOverContexts(ctxs, namespace, options)

			if err != nil {
				return err
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
			return nil
		},
	}

	// Contexts flag
	cmd.Flags().BoolP("favorites", "f", false, "Get all favorited cron jobs")
	// Ordering flags
	cmd.Flags().Bool("by-last-run", false, "Sort chronologically by last run")
	cmd.Flags().Bool("by-last-run-reverse", false, "Sort reverse chronologically by last run")
	cmd.Flags().Bool("by-next-run", false, "Sort cronologically by next scheduled run")
	cmd.Flags().Bool("by-next-run-reverse", false, "Sort reverse chronologically by next scheduled run")

	return cmd
}
