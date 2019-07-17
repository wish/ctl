package parsing

import (
	"errors"
	"github.com/wish/ctl/pkg/client/filter"
	"regexp"
	"strings"
)

// LabelMatch takes in a single label specification string and returns
// a corresponding filter.LabelMatch.
func LabelMatch(s string) (filter.LabelMatch, error) {
	if sub := regexp.MustCompile(`\A(\w+)=(\w+)\z`).FindStringSubmatch(s); len(sub) > 1 {
		return &filter.LabelMatchEq{sub[1], sub[2]}, nil
	} else if sub := regexp.MustCompile(`\A(\w+)!=(\w+)\z`).FindStringSubmatch(s); len(sub) > 1 {
		return &filter.LabelMatchNeq{sub[1], sub[2]}, nil
	} else if sub := regexp.MustCompile(`\A(\w+)\s+in\s*\((\w+(?:,\s*\w+)*)\)\z`).FindStringSubmatch(s); len(sub) > 1 {
		lm := &filter.LabelMatchSetIn{sub[1], nil}
		for _, i := range strings.Split(sub[2], ",") {
			lm.Values = append(lm.Values, strings.TrimSpace(i))
		}
		return lm, nil
	} else { // Did not match any
		return nil, errors.New("no label format found")
	}
}

// LabelMatchSlice is like LabelMatch but handles multiple labels
func LabelMatchSlice(s []string) (filter.LabelMatch, error) {
	if len(s) == 0 {
		return nil, nil
	} else if len(s) == 1 {
		return LabelMatch(s[0])
	}
	labels := make([]filter.LabelMatch, len(s))
	for i, j := range s {
		l, err := LabelMatch(j)
		if err != nil {
			return nil, err
		}
		labels[i] = l
	}
	return &filter.LabelMatchMultiple{labels}, nil
}
