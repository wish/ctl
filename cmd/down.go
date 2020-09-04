package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wish/ctl/pkg/client"
	"os"
	"strings"
)

func downCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down APPNAME",
		Short: "Deletes all ad hoc job with app name (defined in manifest)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user, _ := cmd.Flags().GetString("user")
			appName := args[0]

			// Get hostname to use in job name if not supplied
			if user == "" {
				var err error
				user, err = os.Hostname()
				if err != nil {
					return errors.New("Unable to get hostname of machine")
				}
			}

			// Replace periods with dashes and convert to lower case to follow K8's name constraints
			user = strings.Replace(user, ".", "-", -1)
			user = strings.ToLower(user)

			// Find existing jobs
			job, err := c.FindAdhocJob(appName,user)
			if err != nil {
				return fmt.Errorf("Failed to find jobs: %v", err)
			}

			// Ask the user if they want to delete the current jobs to create a new one
			if job != nil {
				fmt.Printf("Existing job found. Running %s will delete the current jobs, continue? [y/n]\n", appName)
				// Use the prompter from deleteCmd in delete.go
				if !prompter(cmd.InOrStdin()) {
					fmt.Printf("Aborted\n")
					return nil
				}

				fmt.Printf("Deleting %s in %s in context %s...\n", job.Name, job.Namespace, job.Context)
				deleteOptions := client.DeleteOptions{Now: true, DeletionPropagation: true}
				err := c.DeleteJob(job.Context, job.Namespace, job.Name, deleteOptions)
				if err != nil {
						return fmt.Errorf("Error deleting job:\n%v", err)
					}
			}

			fmt.Printf("Job %s deleted.\n", appName)

			return nil
		},
	}
	cmd.Flags().StringP("user", "u", "", "Name that is used for ad hoc jobs. Defaulted to hostname.")

	return cmd
}
