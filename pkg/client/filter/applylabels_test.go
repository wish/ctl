package filter

import "testing"

type testMatchLabelStruct struct {
	match  LabelMatch
	labels map[string]string
	ans    bool
	ans2   bool
}

func (t testMatchLabelStruct) GetLabels() map[string]string {
	return t.labels
}

func TestMatchLabel(t *testing.T) {
	var tests = []testMatchLabelStruct{
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
		{
			nil,
			map[string]string{"a": "b", "b": "b"},
			true,
			true,
		},
		{
			nil,
			nil,
			true,
			true,
		},
		{
			&LabelMatchNeq{"c", "k"},
			nil,
			true,
			true,
		},
		{
			&LabelMatchEq{"c", "k"},
			map[string]string{},
			false,
			true,
		},
	}

	for _, i := range tests {
		if MatchLabel(i, i.match) != i.ans {
			t.Error("Test failed MatchLabel", i.match, i.GetLabels(), i.ans)
		}
		if EmptyOrMatchLabel(i, i.match) != i.ans2 {
			t.Error("Test failed EmptyOrMatchLabel", i.match, i.GetLabels(), i.ans2)
		}
	}
}
