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
	if m == nil {
		return staticlabeled{}
	}
	copy := make(map[string]string)
	for k, v := range m {
		copy[k] = v
	}
	return staticlabeled{labels: copy}
}
