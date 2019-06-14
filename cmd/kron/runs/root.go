package runs

import (
	"github.com/spf13/cobra"
)

func init() {
	// Nothing
}

var RunsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Subcommand on recent runs of a cron job",
	Long: `Operate on the jobs started by a cron job
Has a bunch of subcommand just like kron`,
}
