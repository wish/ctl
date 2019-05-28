package cmd

import (
  //"fmt"
  "github.com/ContextLogic/wishctl/cmd/kron"
)

func init() {
  rootCmd.AddCommand(kron.KronCmd)
}
