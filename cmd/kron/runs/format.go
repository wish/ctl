package runs

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"os"
	"text/tabwriter"
	"gopkg.in/yaml.v2"
)

// TODO: Add better formatting and more fields
func printRunList(lst []client.RunDiscovery) {
	if len(lst) == 0 {
		fmt.Println("No runs found!")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "NAME\tSTATE\tSTART\tEND")

	for _, v := range lst {
		fmt.Fprintf(w, "%s\t", v.Name)
		// State
		if v.Status.Failed > 0 {
			fmt.Fprint(w, "FAILED\t")
		} else if v.Status.CompletionTime != nil {
			fmt.Fprint(w, "SUCCESSFUL\t")
		} else {
			fmt.Fprint(w, "IN PROGRESS\t")
		}
		fmt.Fprintf(w, "%v\t", v.Status.StartTime)
		// END
		if v.Status.CompletionTime != nil {
			fmt.Fprintf(w, "%v\n", v.Status.CompletionTime)
		} else {
			fmt.Fprint(w, "<none>\n")
		}
	}
	w.Flush()
}

func describeRun(run client.RunDiscovery) {
	b, _ := yaml.Marshal(run)
	fmt.Println(string(b))
}
