package filter

func MatchLabel(l Labeled, m LabelMatch) bool {
	if m == nil {
		return true
	}
	return m.Match(l.GetLabels())
}
