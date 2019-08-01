package filter

import (
	"testing"
)

func TestLabelMatchEq(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchEq
		labels map[string]string
		ans    bool
		ans2   bool
	}{
		{
			&LabelMatchEq{"a", "b"},
			map[string]string{"a": "b"},
			true,
			true,
		},
		{
			&LabelMatchEq{"a", "b"},
			map[string]string{"a": "c"},
			false,
			false,
		},
		{
			&LabelMatchEq{"a", "b"},
			map[string]string{"b": "a"},
			false,
			true,
		},
		{
			&LabelMatchEq{"a", "b"},
			nil,
			false,
			true,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed Match", i.match, i.labels, i.ans)
		}
		if i.match.EmptyOrMatch(i.labels) != i.ans2 {
			t.Error("Test failed EmptyOrMatch", i.match, i.labels, i.ans2)
		}
	}
}

func TestLabelMatchNeq(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchNeq
		labels map[string]string
		ans    bool
		ans2   bool
	}{
		{
			&LabelMatchNeq{"a", "b"},
			map[string]string{"a": "b"},
			false,
			false,
		},
		{
			&LabelMatchNeq{"a", "b"},
			map[string]string{"a": "c"},
			true,
			true,
		},
		{
			&LabelMatchNeq{"a", "b"},
			map[string]string{"b": "a"},
			true,
			true,
		},
		{
			&LabelMatchNeq{"a", "b"},
			nil,
			true,
			true,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed Match", i.match, i.labels, i.ans)
		}
		if i.match.EmptyOrMatch(i.labels) != i.ans2 {
			t.Error("Test failed EmptyOrMatch", i.match, i.labels, i.ans2)
		}
	}
}

func TestLabelMatchSetIn(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchSetIn
		labels map[string]string
		ans    bool
		ans2   bool
	}{
		{
			&LabelMatchSetIn{"a", []string{"b"}},
			map[string]string{"a": "b"},
			true,
			true,
		},
		{
			&LabelMatchSetIn{"a", []string{"b", "a"}},
			map[string]string{"a": "c"},
			false,
			false,
		},
		{
			&LabelMatchSetIn{"a", []string{"a", "b", "c"}},
			map[string]string{"a": "b"},
			true,
			true,
		},
		{
			&LabelMatchSetIn{"b", []string{"a", "b", "c"}},
			map[string]string{"a": "b"},
			false,
			true,
		},
		{
			&LabelMatchSetIn{"a", []string{"b"}},
			nil,
			false,
			true,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed Match", i.match, i.labels, i.ans)
		}
		if i.match.EmptyOrMatch(i.labels) != i.ans2 {
			t.Error("Test failed EmptyOrMatch", i.match, i.labels, i.ans2)
		}
	}
}

func TestLabelMatchMultiple(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchMultiple
		labels map[string]string
		ans    bool
		ans2   bool
	}{
		{
			&LabelMatchMultiple{[]LabelMatch{&LabelMatchEq{"a", "b"}, &LabelMatchNeq{"c", "k"}}},
			map[string]string{"a": "b"},
			true,
			true,
		},
		{
			&LabelMatchMultiple{[]LabelMatch{&LabelMatchEq{"a", "b"}, &LabelMatchSetIn{"b", []string{"b", "a"}}}},
			map[string]string{"a": "b", "b": "b"},
			true,
			true,
		},
		{
			&LabelMatchMultiple{[]LabelMatch{&LabelMatchEq{"a", "b"}, &LabelMatchSetIn{"a", []string{"b", "a"}}}},
			map[string]string{"a": "a"},
			false,
			false,
		},
		{
			&LabelMatchMultiple{nil},
			map[string]string{"a": "b", "b": "b"},
			true,
			true,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed Match", i.match, i.labels, i.ans)
		}
		if i.match.EmptyOrMatch(i.labels) != i.ans2 {
			t.Error("Test failed EmptyOrMatch: ", i.match, i.labels, i.ans2)
		}
	}
}
