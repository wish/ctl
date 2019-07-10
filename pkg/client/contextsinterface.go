package client

type contextsGetter interface {
	GetAllContexts() []string
}

type StaticContextsGetter struct {
	contexts []string
}

func (s StaticContextsGetter) GetAllContexts() []string {
	return s.contexts
}
