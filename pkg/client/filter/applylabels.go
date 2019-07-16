package filter

// MatchLabel checks if an object satisfies the requirements of a Labelmatch
func MatchLabel(obj Labeled, lm LabelMatch) bool {
	if lm == nil {
		return true
	}
	return lm.Match(obj.GetLabels())
}
