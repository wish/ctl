package runs

import (
	"fmt"
	"github.com/wish/ctl/pkg/client/types"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/tabwriter"
)

// TODO: Add better formatting and more fields
func printRunList(lst []types.RunDiscovery, labelColumns []string) {
	if len(lst) == 0 {
		fmt.Println("No runs found!")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprint(w, "NAME\tSTATE\tSTART\tEND")
	for _, l := range labelColumns {
		fmt.Fprint(w, "\t", strings.ToUpper(l))
	}
	fmt.Fprintln(w)

	for _, v := range lst {
		fmt.Fprintf(w, "%s", v.Name)
		// State
		if v.Status.Failed > 0 {
			fmt.Fprint(w, "\tFAILED")
		} else if v.Status.CompletionTime != nil {
			fmt.Fprint(w, "\tSUCCESSFUL")
		} else {
			fmt.Fprint(w, "\tIN PROGRESS")
		}
		fmt.Fprintf(w, "\t%v", v.Status.StartTime)
		// END
		if v.Status.CompletionTime != nil {
			fmt.Fprintf(w, "\t%v", v.Status.CompletionTime)
		} else {
			fmt.Fprint(w, "\t<none>")
		}
		// Labels
		for _, l := range labelColumns {
			fmt.Fprint(w, "\t")
			if _, ok := v.Labels[l]; ok {
				fmt.Fprint(w, v.Labels[l])
			}
		}
		fmt.Fprintln(w)
	}
	w.Flush()
}

func describeRun(run types.RunDiscovery) {
	b, _ := yaml.Marshal(run)
	fmt.Println(string(b))
}
