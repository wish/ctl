package filter

import (
	"testing"
)

func TestLabelMatchEq(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchEq
		labels map[string]string
		ans    bool
	}{
		{
			&LabelMatchEq{"a", "b"},
			map[string]string{"a": "b"},
			true,
		},
		{
			&LabelMatchEq{"a", "b"},
			map[string]string{"a": "c"},
			false,
		},
		{
			&LabelMatchEq{"a", "b"},
			map[string]string{"b": "a"},
			false,
		},
		{
			&LabelMatchEq{"a", "b"},
			nil,
			false,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed", i.match, i.labels, i.ans)
		}
	}
}

func TestLabelMatchNeq(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchNeq
		labels map[string]string
		ans    bool
	}{
		{
			&LabelMatchNeq{"a", "b"},
			map[string]string{"a": "b"},
			false,
		},
		{
			&LabelMatchNeq{"a", "b"},
			map[string]string{"a": "c"},
			true,
		},
		{
			&LabelMatchNeq{"a", "b"},
			map[string]string{"b": "a"},
			true,
		},
		{
			&LabelMatchNeq{"a", "b"},
			nil,
			true,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed", i.match, i.labels, i.ans)
		}
	}
}

func TestLabelMatchSetIn(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchSetIn
		labels map[string]string
		ans    bool
	}{
		{
			&LabelMatchSetIn{"a", []string{"b"}},
			map[string]string{"a": "b"},
			true,
		},
		{
			&LabelMatchSetIn{"a", []string{"b", "a"}},
			map[string]string{"a": "c"},
			false,
		},
		{
			&LabelMatchSetIn{"a", []string{"a", "b", "c"}},
			map[string]string{"a": "b"},
			true,
		},
		{
			&LabelMatchSetIn{"b", []string{"a", "b", "c"}},
			map[string]string{"a": "b"},
			false,
		},
		{
			&LabelMatchSetIn{"a", []string{"b"}},
			nil,
			false,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed", i.match, i.labels, i.ans)
		}
	}
}

func TestLabelMatchMultiple(t *testing.T) {
	var tests = []struct {
		match  *LabelMatchMultiple
		labels map[string]string
		ans    bool
	}{
		{
			&LabelMatchMultiple{[]LabelMatch{&LabelMatchEq{"a", "b"}, &LabelMatchNeq{"c", "k"}}},
			map[string]string{"a": "b"},
			true,
		},
		{
			&LabelMatchMultiple{[]LabelMatch{&LabelMatchEq{"a", "b"}, &LabelMatchSetIn{"b", []string{"b", "a"}}}},
			map[string]string{"a": "b", "b": "b"},
			true,
		},
		{
			&LabelMatchMultiple{[]LabelMatch{&LabelMatchEq{"a", "b"}, &LabelMatchSetIn{"a", []string{"b", "a"}}}},
			map[string]string{"a": "a"},
			false,
		},
		{
			&LabelMatchMultiple{nil},
			map[string]string{"a": "b", "b": "b"},
			true,
		},
	}

	for _, i := range tests {
		if i.match.Match(i.labels) != i.ans {
			t.Error("Test failed", i.match, i.labels, i.ans)
		}
	}
}
