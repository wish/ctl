package client

import (
	"encoding/json"
	"fmt"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/pkg/client/types"
)

// FindAdhocJob loops through all valid contexts and returns the first active job found
func (c *Client) FindAdhocJob(appName string, options ListOptions) (*types.JobDiscovery, error) {

	// Get all kubernetes contexts from config file
	config, err := config.GetCtlExt()
	if err != nil {
		return nil, err
	}

	for ctx  := range config {

		if rawruns, ok := config[ctx]["_run"]; ok {
			runs := make(map[string]types.RunDetails)
			err := json.Unmarshal([]byte(rawruns), &runs)
			if err != nil {
				continue
			}

			// Check if the app name exists in the raw runs
			if run, ok := runs[appName]; ok {
				if run.Active {

					// Extract manifest json as struct to parse
					var manifestData types.ManifestDetails
					err = json.Unmarshal([]byte(run.Manifest), &manifestData)
					if err != nil {
						return nil, fmt.Errorf("Error parsing manifestJson: %s", err)
					}

					// Check if a job is already running
					jobs, err := c.ListJobs(ctx, manifestData.Metadata.Namespace, options)
					if err != nil {
						return nil, fmt.Errorf("Failed to search for existing job: %s", err)
					}

					// Return the first job since we limit adhoc pods to one
					if len(jobs) > 0 {
						return &jobs[0], nil
					}
				}
			}
		}
	}
	return nil, nil
}