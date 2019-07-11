package kron

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
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

func GetUnfavoriteCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "unfavorite jobs",
		Short: "Removes job(s) from favorite list",
		Long: `Removes job(s) from favorite list.
	If no jobs are specified, removes selected job.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Behaviour when
			var f map[string]location
			err := viper.UnmarshalKey("favorites", &f)
			if err != nil {
				return err
			}

			for _, job := range args {
				delete(f, job)
			}

			viper.Set("favorites", f)
			viper.WriteConfig()

			return nil
		},
	}
}
