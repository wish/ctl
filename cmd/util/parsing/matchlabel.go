package parsing

import (
	// "fmt"
	"errors"
	"github.com/ContextLogic/ctl/pkg/client/filter"
	"regexp"
	"strings"
)

func LabelMatch(s string) (filter.LabelMatch, error) {
	if sub := regexp.MustCompile(`\A(\w+)=(\w+)\z`).FindStringSubmatch(s); len(sub) > 1 {
		return &filter.LabelMatchEq{sub[1], sub[2]}, nil
	} else if sub := regexp.MustCompile(`\A(\w+)!=(\w+)\z`).FindStringSubmatch(s); len(sub) > 1 {
		return &filter.LabelMatchNeq{sub[1], sub[2]}, nil
	} else if sub := regexp.MustCompile(`\A(\w+)\s+in\s*\(((\w+)(?:,\s*\w+)*)\)\z`).FindStringSubmatch(s); len(sub) > 1 {
		lm := &filter.LabelMatchSetIn{sub[1], nil}
		for _, i := range strings.Split(sub[2], ",") {
			lm.Values = append(lm.Values, strings.TrimSpace(i))
		}
		return lm, nil
	} else { // Did not match any
		return nil, errors.New("No label format found")
	}
}

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
