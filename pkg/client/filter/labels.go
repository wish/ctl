package filter

// Labeled is used to access labels of an object
type Labeled interface {
	GetLabels() map[string]string
}

// LabelMatch is used to filter objects with labels
type LabelMatch interface {
	Match(map[string]string) bool
}

// LabelMatchEq checks if the value of a label is equal to a value
type LabelMatchEq struct {
	Key   string
	Value string
}

// Match checks the label value for equality
func (m *LabelMatchEq) Match(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		return v == m.Value
	}
	return false
}

// LabelMatchNeq checks if the value of a label is non-existent or not equal to a value
type LabelMatchNeq struct {
	Key   string
	Value string
}

// Match checks the label value for non-equality
func (m *LabelMatchNeq) Match(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		return v != m.Value
	}
	return true
}

// LabelMatchSetIn checks if the value of a label is non-existent or is one of multiple values
type LabelMatchSetIn struct {
	Key    string
	Values []string
}

// Match checks the label value belongs in a set
func (m *LabelMatchSetIn) Match(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		for _, vals := range m.Values {
			if v == vals {
				return true
			}
		}
	}
	return false
}

// LabelMatchMultiple checks that multiple label criteria are satisfied
type LabelMatchMultiple struct {
	Matches []LabelMatch
}

// Match checks all nested LabelMatches
func (m *LabelMatchMultiple) Match(labels map[string]string) bool {
	for _, match := range m.Matches {
		if !match.Match(labels) {
			return false
		}
	}
	return true
}
