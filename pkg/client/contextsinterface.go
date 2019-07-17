package client

type contextsGetter interface {
	GetAllContexts() []string
}

// StaticContextsGetter returns a static list of contexts.
// This struct implements the contextsGetter interface
type StaticContextsGetter struct {
	contexts []string
}

// GetAllContexts returns the list of all the clusters
func (s StaticContextsGetter) GetAllContexts() []string {
	return s.contexts
}
