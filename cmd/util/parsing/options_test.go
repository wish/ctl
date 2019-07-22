package parsing

import (
	"github.com/spf13/cobra"
	"testing"
)

func TestAllOptions(t *testing.T) {
	var tests = []struct {
		flags []string
		err   bool
	}{
		{[]string{`-l`, `a=b`}, false},
		{[]string{`-l=a=b`}, false},
		{[]string{`-l`, `a=b`, `-l=a!=b`}, false},
		{[]string{`-l`, `a in (a,b,c)`}, false},
		{[]string{`-l`, `c in(a,z)`, `-l`, `a=b`}, false},
		{[]string{`-l`, `c in (a),b=a`}, false},
		{[]string{`-l`, `a in (x),a`}, true},
		{[]string{`-l`, `123.2.1/a=b`, `-l=a/z!=b`}, false},
		{[]string{`--label`, `a=b`}, false},
		{[]string{`--label`, `a in (a)`}, false},
		{[]string{`--label=a in (a)`}, false},
		{[]string{`--label=a=b`}, false},
		{[]string{`--label=a=b`, `--label=a!=b`}, false},
		{[]string{`--label=a in (a, b, c)`}, false},
		{[]string{`--label=c in(a,z)`, `--label=a=b`}, false},
		{[]string{`--label=c in (a)`}, false},
		{[]string{`--label=a/c in (a)`}, false},
	}

	for _, test := range tests {
		cmd := &cobra.Command{
			Args: cobra.ExactArgs(0),
			RunE: func(cmd *cobra.Command, args []string) error {
				_, err1 := ListOptions(cmd)
				_, err2 := GetOptions(cmd)
				_, err3 := LogOptions(cmd)

				if err1 != nil {
					return err1
				}
				if err2 != nil {
					return err1
				}
				if err3 != nil {
					return err1
				}
				return nil
			},
		}
		cmd.Flags().StringArrayP("label", "l", nil, "blah")
		cmd.SetArgs(test.flags)

		_, err := cmd.ExecuteC()

		if err != nil {
			if !test.err {
				t.Error("Encountered unexpected error on", test.flags, ":", err)
			}
		} else if test.err {
			t.Error("Expected an error but did not encounter one.", test.flags)
		}
	}
}
