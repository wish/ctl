package kron

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wish/ctl/pkg/client"
)

func favoriteCmd(c *client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "favorite [jobs] [flags]",
		Short: "Adds jobs to favorite list",
		Long: `Adds specified job(s) to the favorite list. If no job was specified the selected job is added.
A namespace and contexts can be specified to limit matches.`,
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
			// args/flags
			ctxs, _ := cmd.Flags().GetStringSlice("context")
			nss, _ := cmd.Flags().GetString("namespace")

			f, err := getFavorites()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				selected, err := getSelected()
				if err != nil {
					return err
				}
				if l, ok := f[selected.Name]; ok {
					cmd.Println(overrideFavoriteMessage(selected.Name, l))
				}
				f[selected.Name] = selected.Location
			} else {
				for _, job := range args {
					if l, ok := f[job]; ok {
						cmd.Println(overrideFavoriteMessage(job, l))
					}
					f[job] = location{ctxs, nss}
				}
			}

			viper.Set("favorites", f)
			viper.WriteConfig()

			return nil
		},
	}
}
