package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/clusterext"
)

type runDetails struct {
	Resources    resource `json:"resources"`
	Active       bool     `json:"active"`
	Manifest     string   `json:"manifest"`
	PreLogin	 [][]string `json:"pre_login_command,omitempty"`
	LoginCommand []string `json:"login_command"`
}

type resource struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

const (
	// DefaultDeadline - default amount (in secs) of time to keep the ad hoc pods running
	DefaultDeadline string = "43200" // 12 hours = 60 * 60 * 12
	// MaxDeadline sets the max deadline for a job (defaulted to 1 day)
	MaxDeadline int = 60 * 60 * 24
	// DefaultCPU is used to template cpu in manifest file if no default or flag is found
	DefaultCPU string = "0.5"
	// DefaultMemory is used to template memory in manifest file if no default or flag is found
	DefaultMemory string = "128Mi"
)

func upCmd(c *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up APPNAME [flags]",
		Short: "Creates an ad hoc job with app name (defined in manifest)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			labelMatch, _ := parsing.LabelMatchFromCmd(cmd)
			deadline, _ := cmd.Flags().GetString("deadline")
			cpu, _ := cmd.Flags().GetString("cpu")
			memory, _ := cmd.Flags().GetString("memory")
			user, _ := cmd.Flags().GetString("user")

			// Check for valid input for deadline and set default if needed
			if deadlineInt, err := strconv.Atoi(deadline); err != nil || deadlineInt < 1 || deadlineInt > MaxDeadline {
				deadline = DefaultDeadline
			}

			appName := args[0]

			// Get all kubernetes contexts from config file
			m, err := config.GetCtlExt()
			if err != nil {
				return err
			}
			e := clusterext.Extension{ClusterExt: m, K8Envs: nil}
			ctxs := e.GetFilteredContexts(labelMatch)

			// Random shuffle contexts
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(ctxs), func(i, j int) { ctxs[i], ctxs[j] = ctxs[j], ctxs[i] })

			for _, ctx := range ctxs {

				if rawruns, ok := m[ctx]["_run"]; ok {
					runs := make(map[string]runDetails)
					err := json.Unmarshal([]byte(rawruns), &runs)
					if err != nil { // bad
						continue
					}

					// Get a list of available app names
					availableAppNames := make([]string, len(runs))
					i := 0
					for k := range runs {
						availableAppNames[i] = k
						i++
					}

					// Check if the app name exists in the raw runs
					if run, ok := runs[appName]; ok {
						if run.Active {
							// Get hostname to use in job name if not supplied
							if user == "" {
								user, err = os.Hostname()
								if err != nil {
									return errors.New("Unable to get hostname of machine")
								}
							}

							//Replace periods with dashes to follow K8's name constraints
							user = strings.Replace(user, ".", "-", -1) 
							
							// First, let's check if a job is already running. We want to limit 1 job per user.
							// We find the job using its name '<app name>-<host name>' eg 'foo-bar'
							jobs, err := c.FindJobs([]string{}, "", []string{fmt.Sprintf("%s-%s", appName, user)},
								client.ListOptions{},
							)
							if err != nil {
								return fmt.Errorf("Failed to find jobs: %v", err)
							}

							// Ask the user if they want to delete the current jobs to create a new one
							if len(jobs) > 0 {
								fmt.Printf("\nExisting jobs: (%d) found. Running %s will delete the current jobs, continue? [y/n]\n",
									len(jobs), appName)
								// Use the prompter from deleteCmd in delete.go
								if !prompter(cmd.InOrStdin()) {
									fmt.Printf("Aborted\n")
									return nil
								}

								for _, job := range jobs {
									deleteOptions := client.DeleteOptions{Now: true, DeletionPropagation: true}
									c.DeleteJob(job.Context, job.Namespace, job.Name, deleteOptions)
								}
							}

							// Template out hostname into job name in manifest
							manifest := regexp.MustCompile(`({USER})`).ReplaceAllString(run.Manifest, user)
							// Template active deadline seconds into manifest
							manifest = regexp.MustCompile(`("{ACTIVE_DEADLINE_SECONDS}")`).ReplaceAllString(manifest, deadline)
							// Template out resources into the manifest
							if cpu != "" {
								run.Resources.CPU = cpu
							} else if run.Resources.CPU == "" {
								run.Resources.CPU = DefaultCPU
							}
							manifest = regexp.MustCompile(`({CPU})`).ReplaceAllString(manifest, run.Resources.CPU)
							if memory != "" {
								run.Resources.Memory = memory
							} else if run.Resources.Memory == "" {
								run.Resources.Memory = DefaultMemory
							}
							manifest = regexp.MustCompile(`({MEMORY})`).ReplaceAllString(manifest, run.Resources.Memory)

							// Add context flag in case the namespace does not exist in current cluster
							context := fmt.Sprintf("--context=%s", ctx)

							// Pass the manifest into a reader for stdin
							r := strings.NewReader(manifest)
							command := exec.Command("kubectl", "apply", "-f", "-", context)
							command.Stdout = os.Stdout
							command.Stderr = os.Stderr
							command.Stdin = r

							fmt.Printf("Running %v with a deadline of %v seconds. CPU: %v, Memory: %v...\n",
								appName, deadline, run.Resources.CPU, run.Resources.Memory)
							fmt.Printf("JOB: %v\nCONTEXT: %v\n\nUse `ctl login %s` to sh into your pod\n\n",
								appName+"-"+user, ctx, appName)

							return command.Run()
						}
						// If the app name is not active, let's just print its manifest
						cmd.Printf("WARN: App name not active in %v. Manifest file: %v\n", ctx, run.Manifest)
					} else {
						// App name not found, list out available app names to the user
						fmt.Printf("Available app names: %v\n", availableAppNames)
						return errors.New("App name not found")
					}
				}
			}
			return errors.New("no command found to run or not authorized to access cluster")
		},
	}

	cmd.Flags().String("deadline", "", "Time pod will stay alive in seconds")
	cmd.Flags().String("cpu", "", "CPU for pod, default is "+DefaultCPU+". eg. --cpu=0.5")
	cmd.Flags().String("memory", "", "Memory for pod, default is "+DefaultMemory+". eg --memory=4.0Gi")
	cmd.Flags().StringP("user", "u", "", "Name that is used for ad hoc jobs. Defaulted to hostname.")

	return cmd
}
