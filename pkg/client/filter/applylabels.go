package filter

// MatchLabel checks if an object satisfies the requirements of a LabelMatch
func MatchLabel(obj Labeled, lm LabelMatch) bool {
	if lm == nil {
		return true
	}
	return lm.Match(obj.GetLabels())
}

// EmptyOrMatchLabel checks EmptyOrMatch on a LabelMatch
func EmptyOrMatchLabel(obj Labeled, lm LabelMatch) bool {
	if lm == nil {
		return true
	}
	return lm.EmptyOrMatch(obj.GetLabels())
}
