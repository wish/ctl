package kron

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func init() {
	KronCmd.AddCommand(favoriteCmd)
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

var favoriteCmd = &cobra.Command{
	Use:   "favorite [jobs] [flags]",
	Short: "Adds jobs to favorite list",
	Long: `Adds specified job(s) to the favorite list. If no job was specified the selected job is added.
A namespace and contexts can be specified to limit matches.`,
	Run: func(cmd *cobra.Command, args []string) {
		// args/flags
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		nss, _ := cmd.Flags().GetString("namespace")

		f, err := getFavorites()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if len(args) == 0 {
			selected, err := getSelected()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if l, ok := f[selected.Name]; ok {
				fmt.Println(overrideFavoriteMessage(selected.Name, l))
			}
			f[selected.Name] = selected.Location
		} else {
			for _, job := range args {
				if l, ok := f[job]; ok {
					fmt.Println(overrideFavoriteMessage(job, l))
				}
				f[job] = location{ctxs, nss}
			}
		}

		viper.Set("favorites", f)
		viper.WriteConfig()
	},
}
