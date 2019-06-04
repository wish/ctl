package kron

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	KronCmd.AddCommand(unfavoriteCmd)
	viper.SetDefault("favorites", make(map[string]location))
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.kron")
	err := viper.ReadInConfig()
	if err != nil {
		// Write config file
		fmt.Println("Creating new config file")
		createConfig()
		// panic(err.Error())
	}
}

var unfavoriteCmd = &cobra.Command{
	Use:   "unfavorite job...",
	Short: "Removes job(s) from favorite list",
	Long:  "Removes job(s) from favorite list. If no jobs are specified, removes selected job. If job is selected, opens a list to choose to remove from.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// args/flags

		// Behaviour when
		var f map[string]location
		err := viper.UnmarshalKey("favorites", &f)
		if err != nil {
			fmt.Println(err.Error())
		}

		for _, job := range args {
			delete(f, job)
		}

		viper.Set("favorites", f)
		viper.WriteConfig()
	},
}
