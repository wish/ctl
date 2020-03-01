package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
)

func downCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down APPNAME",
		Short: "Deletes all ad hoc job with app name (defined in manifest)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]

			// Get hostname to use in job name
			user, err := os.Hostname()
			if err != nil {
				return errors.New("Unable to get hostname of machine")
			}

			// Find existing jobs
			jobs, err := c.FindJobs([]string{}, "", []string{fmt.Sprintf("%s-%s", appName, user)},
				client.ListOptions{},
			)

			// Ask the user if they want to delete the current jobs to create a new one
			if len(jobs) > 0 {
				fmt.Printf("Existing jobs: (%d) found. Running %s will delete the current jobs, continue? [y/n]\n",
					len(jobs), appName)
				// Use the prompter from deleteCmd in delete.go
				if !prompter(cmd.InOrStdin()) {
					fmt.Printf("Aborted\n")
					return nil
				}

				for _, job := range jobs {
					fmt.Printf("Deleting %s in %s in context %s...\n", job.Name, job.Namespace, job.Context)
					deleteOptions := client.DeleteOptions{Now: true, DeletionPropagation: true}
					err := c.DeleteJob(job.Context, job.Namespace, job.Name, deleteOptions)
					if err != nil {
						return fmt.Errorf("Error deleting job:\n%v", err)
					}
				}
			}

			fmt.Printf("All %s jobs deleted.\n", appName)

			return nil
		},
	}

	return cmd
}
