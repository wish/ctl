package kron

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	favoriteCmd.Flags().StringSliceP("contexts", "c", []string{}, "Specific contexts to list cronjobs from")
	favoriteCmd.Flags().StringSliceP("namespaces", "n", []string{}, "Specific namespaces to list cronjobs from within contexts")
}

var favoriteCmd = &cobra.Command{
	Use:   "favorite job",
	Short: "Adds a job to favorite list",
	Long:  "Adds specified job(s) to the favorite list. If no job was specified the selected job is added.",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// args/flags
		job := args[0]
		ctxs, _ := cmd.Flags().GetStringSlice("contexts")
		nss, _ := cmd.Flags().GetStringSlice("namespaces")

		// Behaviour when
		var f map[string]location
		err := viper.UnmarshalKey("favorites", &f)
		if err != nil {
			fmt.Println(err.Error())
		}

		if l, ok := f[job]; ok {
			fmt.Printf("Job \"%s\" is already in favorites with:\n", job)
			fmt.Printf("Contexts: %v\n", l.Contexts)
			fmt.Printf("Namespaces: %v\n", l.Namespaces)
			fmt.Println("Overriding entry...")
		}
		f[job] = location{ctxs, nss}

		viper.Set("favorites", f)
		viper.WriteConfig()
	},
}
