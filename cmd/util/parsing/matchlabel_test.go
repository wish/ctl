package parsing

import (
	"github.com/wish/ctl/pkg/client/filter"
	"testing"
)

func TestLabelMatchSlice(t *testing.T) {
	var labeltests = []struct {
		str []string
		cat string // err, eq, neq, setin or multiple
	}{
		{[]string{"abc=abc"}, "eq"},
		{[]string{"pft"}, "err"},
		{[]string{"a=b", "q!=x"}, "multiple"},
		{[]string{"="}, "err"},
		{[]string{"label!=x"}, "neq"},
		{[]string{"s in (1, 2, 3)"}, "setin"},
		{[]string{"s in(a)"}, "setin"},
	}

	for _, i := range labeltests {
		lm, err := LabelMatchSlice(i.str)
		if err != nil {
			if i.cat != "err" {
				t.Error(err.Error())
			}
			continue
		}

		switch lm.(type) {
		case *filter.LabelMatchEq:
			if i.cat != "eq" {
				t.Error("Parsed label match not expected!")
			}
		case *filter.LabelMatchNeq:
			if i.cat != "neq" {
				t.Error("Parsed label match not expected!")
			}
		case *filter.LabelMatchSetIn:
			if i.cat != "setin" {
				t.Error("Parsed label match not expected!")
			}
		case *filter.LabelMatchMultiple:
			if i.cat != "multiple" {
				t.Error("Parsed label match not expected!")
			}
		}
	}
}
