package kron

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	KronCmd.AddCommand(selectCmd)
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

var selectCmd = &cobra.Command{
	Use:   "select job [flags]",
	Short: "Uses list to select a job to operate on",
	Long: `Uses list to select a job on which other commands can conveniently operate on.
A namespace and contexts can be specified to limit matches.
If namespace/contexts are not specified, usage will match with all results.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// args/flags
		job := args[0]
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		namespace, _ := cmd.Flags().GetString("namespace")

		var s selectedJob
		err := viper.UnmarshalKey("selected", &s)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		viper.Set("selected", selectedJob{job, location{ctxs, namespace}})
		viper.WriteConfig()
	},
}
