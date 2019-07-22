package filter

// LabelMatch is used to filter objects with labels
type LabelMatch interface {
	Match(map[string]string) bool
	EmptyOrMatch(map[string]string) bool
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

// EmptyOrMatch checks if the label can be filtered on existent fields
func (m *LabelMatchEq) EmptyOrMatch(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		return v == m.Value
	}
	return true
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

// EmptyOrMatch checks if the label can be filtered on existent fields
func (m *LabelMatchNeq) EmptyOrMatch(labels map[string]string) bool {
	return m.Match(labels)
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

// EmptyOrMatch checks if the label can be filtered on existent fields
func (m *LabelMatchSetIn) EmptyOrMatch(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		for _, vals := range m.Values {
			if v == vals {
				return true
			}
		}
		return false
	}
	return true
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

// EmptyOrMatch checks if the label can be filtered on existent fields
func (m *LabelMatchMultiple) EmptyOrMatch(labels map[string]string) bool {
	for _, match := range m.Matches {
		if !match.EmptyOrMatch(labels) {
			return false
		}
	}
	return true
}
