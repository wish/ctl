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
	if sub := regexp.MustCompile(`\A((?:[\w-.]+/)?[\w-.]+)\s*=\s*([\w-.]+)\z`).FindStringSubmatch(s); len(sub) > 1 {
		return &filter.LabelMatchEq{sub[1], sub[2]}, nil
	} else if sub := regexp.MustCompile(`\A((?:[\w-.]+/)?[\w-.]+)\s*!=\s*([\w-.]+)\z`).FindStringSubmatch(s); len(sub) > 1 {
		return &filter.LabelMatchNeq{sub[1], sub[2]}, nil
	} else if sub := regexp.MustCompile(`\A((?:[\w-.]+/)?[\w-.]+)\s+in\s*\(([\w-.]+(?:,\s*[\w-.]+)*)\)\z`).FindStringSubmatch(s); len(sub) > 1 {
		lm := &filter.LabelMatchSetIn{sub[1], nil}
		for _, i := range strings.Split(sub[2], ",") {
			lm.Values = append(lm.Values, strings.TrimSpace(i))
		}
		return lm, nil
	} else { // Did not match any
		return nil, errors.New(`no label format found for "` + s + `"`)
	}
}

const allLabels string = `(?:[\w-.]+/)?[\w-.]+\s*=\s*[\w-.]+|(?:[\w-.]+/)?[\w-.]+\s*!=\s*[\w-.]+|(?:[\w-.]+/)?[\w-.]+\s+in\s*\([\w-.]+(?:,\s*[\w-.]+)*\)`
const allCat string = `\A` + allLabels + `(,` + allLabels + `)*\z`

// LabelMatchSlice is like LabelMatch but handles multiple labels
func LabelMatchSlice(values []string) (filter.LabelMatch, error) {
	var split []string
	for _, s := range values {
		if !regexp.MustCompile(allCat).MatchString(s) {
			return nil, errors.New(`Label "` + s + `" is not of proper format"`)
		}
		for _, i := range regexp.MustCompile(allLabels).FindAllStringSubmatch(s, -1) {
			if len(i) == 1 {
				split = append(split, i[0])
			}
		}
	}
	if len(split) == 0 {
		return nil, nil
	} else if len(split) == 1 {
		return LabelMatch(split[0])
	}
	labels := make([]filter.LabelMatch, len(split))
	for i, j := range split {
		l, err := LabelMatch(j)
		if err != nil {
			return nil, err
		}
		labels[i] = l
	}
	return &filter.LabelMatchMultiple{labels}, nil
}
