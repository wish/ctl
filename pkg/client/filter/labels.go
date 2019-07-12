package filter

type Labeled interface {
	GetLabels() map[string]string
}

type LabelMatch interface {
	Match(map[string]string) bool
}

type LabelMatchEq struct {
	Key   string
	Value string
}

func (m *LabelMatchEq) Match(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		return v == m.Value
	}
	return false
}

type LabelMatchNeq struct {
	Key   string
	Value string
}

func (m *LabelMatchNeq) Match(labels map[string]string) bool {
	if v, ok := labels[m.Key]; ok {
		return v != m.Value
	}
	return true
}

type LabelMatchSetIn struct {
	Key    string
	Values []string
}

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

type LabelMatchMultiple struct {
	Matches []LabelMatch
}

func (m *LabelMatchMultiple) Match(labels map[string]string) bool {
	for _, match := range m.Matches {
		if !match.Match(labels) {
			return false
		}
	}
	return true
}
