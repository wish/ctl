package filter

// Labeled is used to access labels of an object
type Labeled interface {
	GetLabels() map[string]string
}

// Temporary objects
type staticlabeled struct {
	labels map[string]string
}

func (s staticlabeled) GetLabels() map[string]string {
	return s.labels
}

// GetLabeled returns a Labeled object that returns the passed map
func GetLabeled(m map[string]string) Labeled {
	return staticlabeled{labels: m}
}
