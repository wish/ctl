package labelforger

import (
	"github.com/wish/ctl/pkg/client/types"
)

// LabelForger is an object that inserts cluster entries into the labels of objects
type LabelForger struct {
	ClusterExt map[string]map[string]string
}

func transformObjectMeta(m map[string]string, labels *map[string]string) {
	if m == nil {
		return
	}
	if *labels == nil {
		*labels = make(map[string]string)
	}
	for k, v := range m {
		if _, ok := (*labels)[k]; !ok {
			(*labels)[k] = v
		}
	}
}

// Transform inserts any labels from the cluster
func (l LabelForger) Transform(i interface{}) {
	if len(l.ClusterExt) == 0 {
		return
	}
	switch v := i.(type) {
	case *types.CronJobDiscovery:
		transformObjectMeta(l.ClusterExt[v.Context], &(v.Labels))
	case *types.PodDiscovery:
		transformObjectMeta(l.ClusterExt[v.Context], &(v.Labels))
	case *types.RunDiscovery:
		transformObjectMeta(l.ClusterExt[v.Context], &(v.Labels))
	default:
		panic("Unknown object")
	}
}
