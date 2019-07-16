package kron

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client"
)

func unfavoriteCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "unfavorite jobs",
		Short: "Removes job(s) from favorite list",
		Long: `Removes job(s) from favorite list.
If no jobs are specified, removes selected job.`,
		Args: cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.SetDefault("favorites", make(map[string]location))
			viper.SetConfigName("config")
			viper.AddConfigPath("$HOME/.kron")
			err := viper.ReadInConfig()
			if err != nil {
				// Write config file
				cmd.Println("Creating new config file")
				createConfig()
			}
		},
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
