package clusterext

import (
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	"log"
)

// Extension is an object that inserts cluster entries into the labels of objects
// it also supports filtering clusters by labels
type Extension struct {
	ClusterExt map[string]map[string]string
}

// EmptyExtension returns an Extension object that works as if no extensions are set
func EmptyExtension(clusters []string) Extension {
	m := make(map[string]map[string]string)
	for _, c := range clusters {
		m[c] = nil
	}
	return Extension{ClusterExt: m}
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
func (e Extension) Transform(i interface{}) {
	if len(e.ClusterExt) == 0 {
		return
	}
	switch v := i.(type) {
	case *types.CronJobDiscovery:
		transformObjectMeta(e.ClusterExt[v.Context], &(v.Labels))
	case *types.PodDiscovery:
		transformObjectMeta(e.ClusterExt[v.Context], &(v.Labels))
	case *types.JobDiscovery:
		transformObjectMeta(e.ClusterExt[v.Context], &(v.Labels))
	case *types.ConfigMapDiscovery:
		transformObjectMeta(e.ClusterExt[v.Context], &(v.Labels))
	case *types.DeploymentDiscovery:
		transformObjectMeta(e.ClusterExt[v.Context], &(v.Labels))
	case *types.ReplicaSetDiscovery:
		transformObjectMeta(e.ClusterExt[v.Context], &(v.Labels))
	default:
		log.Printf("unsupported type %T\n", v)
	}
}

// GetFilteredContexts returns all the contexts that may match the label
func (e Extension) GetFilteredContexts(l filter.LabelMatch) []string {
	contexts := make([]string, 0)
	for c, m := range e.ClusterExt {
		if v, ok := m["hidden"]; (!ok || v != "true") && filter.EmptyOrMatchLabel(filter.GetLabeled(m), l) {
			contexts = append(contexts, c)
		}
	}
	return contexts
}

// FilterContexts further filters contexts
func (e Extension) FilterContexts(contexts []string, l filter.LabelMatch) []string {
	ret := make([]string, 0)
	for _, ctx := range contexts {
		if m, ok := e.ClusterExt[ctx]; ok {
			if v, ok := m["hidden"]; (!ok || v != "true") && filter.EmptyOrMatchLabel(filter.GetLabeled(m), l) {
				ret = append(ret, ctx)
			}
		} else { // Can't find, add to list anyways
			ret = append(ret, ctx)
		}
	}
	return ret
}
