package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("selected", make(map[string]location))
	viper.AddConfigPath("$HOME/.kron")
	err := viper.ReadInConfig()
	if err != nil {
		// Write config file
		fmt.Println("Creating new config file")
		createConfig()
		// panic(err.Error())
	}
}

func GetSelectCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "select job [flags]",
		Short: "Uses list to select a job to operate on",
		Long: `Uses list to select a job on which other commands can conveniently operate on.
	A namespace and contexts can be specified to limit matches.
	If namespace/contexts are not specified, usage will match with all results.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// args/flags
			job := args[0]
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			namespace, _ := cmd.Flags().GetString("namespace")

			var s selectedJob
			err := viper.UnmarshalKey("selected", &s)
			if err != nil {
				return err
			}

			viper.Set("selected", selectedJob{job, location{ctxs, namespace}})
			viper.WriteConfig()

			return nil
		},
	}
}
