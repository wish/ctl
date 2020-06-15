package parsing

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/filter"
	v1 "k8s.io/api/core/v1"
	"regexp"
	"strings"
)

// LabelMatchFromCmd automatically parses the "label" flag from a command
// and returns the filtering.LabelMatch specified.
func LabelMatchFromCmd(cmd *cobra.Command) (filter.LabelMatch, error) {
	s, _ := cmd.Flags().GetStringArray("label")
	for _, label := range viper.GetStringSlice("label_flags") {
		if v, err := cmd.Flags().GetString(label); err == nil && len(v) > 0 {
			s = append(s, label+"="+v)
		}
	}
	return LabelMatchSlice(s)
}

//StatusFromCmd automatically parses the status flag from a command
//and returns the filtering.StatusMatch specified.
func StatusFromCmd(cmd *cobra.Command) (filter.StatusMatch, error) {
	s , err := cmd.Flags().GetString("status")
	status := filter.StatusMatch{State: v1.PodPhase(s)}
	return status, err
}

// ListOptions parses a client.ListOptions from a command
func ListOptions(cmd *cobra.Command, searches []string) (client.ListOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	if err != nil {
		return client.ListOptions{}, err
	}

	s, err := StatusFromCmd(cmd)
	if err != nil {
		return client.ListOptions{}, err
	}

	var re *regexp.Regexp
	if len(searches) > 0 {
		// Check that each individual search is a valid regex
		for _, s := range searches {
			if _, err := regexp.Compile(s); err != nil {
				return client.ListOptions{}, err
			}
		}
		re, err = regexp.Compile("^(" + strings.Join(searches, ")|^(") + ")")
	}
	if err != nil {
		return client.ListOptions{}, err
	}
	return client.ListOptions{LabelMatch: l, StatusMatch: s, Search: re}, nil
}

// GetOptions parses a client.GetOptions from a command
func GetOptions(cmd *cobra.Command) (client.GetOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	return client.GetOptions{LabelMatch: l}, err
}

// LogOptions parses a client.LogOptions from a command
func LogOptions(cmd *cobra.Command, searches []string) (client.LogOptions, error) {
	l, err := LabelMatchFromCmd(cmd)
	var re *regexp.Regexp
	if len(searches) > 0 {
		// Check that each individual search is a valid regex
		for _, s := range searches {
			if _, err := regexp.Compile(s); err != nil {
				return client.LogOptions{}, err
			}
		}
		re, err = regexp.Compile("^(" + strings.Join(searches, ")|^(") + ")")
	}
	if err != nil {
		return client.LogOptions{}, err
	}
	return client.LogOptions{LabelMatch: l, Follow: false, Search: re}, nil
}
