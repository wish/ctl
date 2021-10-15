package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wish/ctl/pkg/client/types"
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

const (
	// DefaultDeadline - default amount (in secs) of time to keep the ad hoc pods running
	DefaultDeadline string = "43200" // 12 hours = 60 * 60 * 12
	// MaxDeadline sets the max deadline for a job (defaulted to 2 days)
	MaxDeadline int = 60 * 60 * 24 * 2
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
			image, _ := cmd.Flags().GetString("image")
			container, _ := cmd.Flags().GetString("container")
			memory, _ := cmd.Flags().GetString("memory")
			user, _ := cmd.Flags().GetString("user")
			context, _ := cmd.Flags().GetStringSlice("context")

			// Check for valid input for deadline and set default if needed
			if deadlineInt, err := strconv.Atoi(deadline); err != nil || deadlineInt < 1 || deadlineInt > MaxDeadline {
				deadline = DefaultDeadline
			}

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

			appName := args[0]

			// Get all kubernetes contexts from config file
			m, err := config.GetCtlExt()
			if err != nil {
				return err
			}
			e := clusterext.Extension{ClusterExt: m, K8Envs: nil}

			var ctxs []string
			if len(context) > 0 {
				ctxs = context
			} else {
				ctxs = e.GetFilteredContexts(labelMatch)
			}

			// Random shuffle contexts
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(ctxs), func(i, j int) { ctxs[i], ctxs[j] = ctxs[j], ctxs[i] })

			for _, ctx := range ctxs {

				if rawruns, ok := m[ctx]["_run"]; ok {
					runs := make(map[string]types.RunDetails)
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
							// Extract manifest json as struct to parse
							var manifestData types.ManifestDetails
							err = json.Unmarshal([]byte(run.Manifest), &manifestData)
							if err != nil {
								return fmt.Errorf("Error parsing manifestJson: %s", err)
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

							// Check if a job is already running. We want to limit 1 job per user.
							// Get the namespace for the job from the manifest
							// We find the job using its name '<app name>-<host name>' eg 'foo-bar'
							jobs, err := c.FindJobs([]string{}, manifestData.Metadata.Namespace, []string{fmt.Sprintf("%s-%s", appName, user)},
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
									err := c.DeleteJob(job.Context, job.Namespace, job.Name, deleteOptions)
									if err != nil {
										fmt.Printf("Failed to delete job: %v", err)
									}
								}
							}

							// Update image for specified container if user wants to override default images
							if image != "" {
								if container == "" {
									return fmt.Errorf("Container name missing. To update image of a container, `container` and `image` flags are required")
								}

								for i, containerDetail := range manifestData.Spec.Template.Spec.Containers {
									if containerDetail.Name == container {
										manifest = strings.ReplaceAll(manifest, containerDetail.Image, image)
										break;
									} else if i == len(manifestData.Spec.Template.Spec.Containers)-1 { // If container is not found
										return fmt.Errorf("Container %s not found", container)
									}
								}
							}

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

	cmd.Flags().String("deadline", "", "Time pod will stay alive in seconds. Defaulted to 43200 seconds (12 hours)")
	cmd.Flags().String("cpu", "", "CPU for pod, default is "+DefaultCPU+". eg. --cpu=0.5")
	cmd.Flags().String("memory", "", "Memory for pod, default is "+DefaultMemory+". eg --memory=4.0Gi")
	cmd.Flags().StringP("user", "u", "", "Name that is used for ad hoc jobs. Defaulted to hostname.")
	cmd.Flags().StringP("image", "i", "", "Custom image to launch the pod with. Specify the container that is to be updated with image using `container` flag ")
	cmd.Flags().StringP("container", "c", "", "Name of container that is to be updated with custom image")

	return cmd
}
