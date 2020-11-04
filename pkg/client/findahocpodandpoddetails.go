package client

import (
	"encoding/json"
	"fmt"
	"github.com/wish/ctl/cmd/util/config"
	"github.com/wish/ctl/pkg/client/types"
)

// FindAdhocPodAndAppDetails loops through all valid contexts and returns the first adhoc pod found along with context details about the pod
func (c *Client) FindAdhocPodAndAppDetails(appName string, options ListOptions) (*types.PodDiscovery, *types.ManifestDetails, *types.RunDetails, error) {

	// Get all kubernetes contexts from config file
	config, err := config.GetCtlExt()
	if err != nil {
		return nil, nil, nil, err
	}

	for ctx := range config {

		if rawruns, ok := config[ctx]["_run"]; ok {
			runs := make(map[string]types.RunDetails)
			err := json.Unmarshal([]byte(rawruns), &runs)
			if err != nil {
				continue
			}

			// Check if appName is valid
			if run, ok := runs[appName]; ok {
				if run.Active {

					// Extract manifest json as struct to parse
					var manifestData types.ManifestDetails
					err = json.Unmarshal([]byte(run.Manifest), &manifestData)
					if err != nil {
						return nil, nil, nil, fmt.Errorf("Error parsing manifestJson: %s", err)
					}

					// Check if a job is already running
					pods, err := c.ListPods(ctx, manifestData.Metadata.Namespace, options)
					if err != nil {
						return nil, nil, nil, fmt.Errorf("Failed to search for existing job: %s", err)
					}

					// Return the first pod since we limit adhoc pods to one
					if len(pods) > 0 {
						return &pods[0], &manifestData, &run, nil
					}
				}
			}
		}
	}
	return nil, nil, nil, nil
}