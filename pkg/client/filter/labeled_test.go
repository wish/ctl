package filter

import (
	"reflect"
	"testing"
)

func TestLabeled(t *testing.T) {
	tests := []map[string]string{
		nil,
		{},
		{"a": "b"},
		{"a": "b", "c": "d"},
	}

	for _, test := range tests {
		l := GetLabeled(test)
		if !reflect.DeepEqual(test, l.GetLabels()) {
			t.Error("GetLabeled failed on", test)
		}
	}
	// Custom test for deep copy.
	a := map[string]string{"a": "b"}
	l := GetLabeled(a)
	b := a
	a = map[string]string{"a": "c"}
	if reflect.DeepEqual(a, l.GetLabels()) || !reflect.DeepEqual(b, l.GetLabels()) {
		t.Error("GetLabeled failed on custom dynamic test")
		// If this test is failing, then GetLabeled made a shallow copy of the map
	}
}
