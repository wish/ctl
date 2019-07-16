package filter

import "testing"

type testMatchLabelStruct struct {
	match  *LabelMatchMultiple
	labels map[string]string
	ans    bool
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
		if MatchLabel(i, i.match) != i.ans {
			t.Error("Test failed", i.match, i.GetLabels(), i.ans)
		}
	}
}
