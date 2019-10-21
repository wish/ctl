package clusterext

import (
	"github.com/wish/ctl/pkg/client/filter"
	"github.com/wish/ctl/pkg/client/types"
	batchv1 "k8s.io/api/batch/v1"
	batchv1b1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sort"
	"testing"
)

// Helper
func createObj(context, name, namespace, res string, labels map[string]string) interface{} {
	objectmeta := metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
	switch res {
	case "pod":
		return &types.PodDiscovery{
			Context: context,
			Pod: corev1.Pod{
				ObjectMeta: objectmeta,
			},
		}
	case "cronjob":
		return &types.CronJobDiscovery{
			Context: context,
			CronJob: batchv1b1.CronJob{
				ObjectMeta: objectmeta,
			},
		}
	case "job":
		return &types.JobDiscovery{
			Context: context,
			Job: batchv1.Job{
				ObjectMeta: objectmeta,
			},
		}
	}
	return nil
}

func TestTransform(t *testing.T) {
	var tests = []struct {
		lf   map[string]map[string]string
		objs []struct {
			obj interface{}
			ans map[string]string
		}
	}{
		{
			lf: nil,
			objs: []struct {
				obj interface{}
				ans map[string]string
			}{
				{createObj("anything", "name", "default", "pod", nil), nil},
			},
		},
		{
			lf: map[string]map[string]string{
				"c1": map[string]string{
					"foo":   "bar",
					"hello": "world",
				},
				"c2": map[string]string{
					"pft": "tfp",
				},
				"c3": nil,
			},
			objs: []struct {
				obj interface{}
				ans map[string]string
			}{
				{createObj("c2", "name", "default", "cronjob", nil), map[string]string{"pft": "tfp"}},
				{createObj("c1", "name", "default", "job", map[string]string{"a": "b"}), map[string]string{"a": "b", "foo": "bar", "hello": "world"}},
				{createObj("c3", "wow", "bad", "pod", nil), nil},
			},
		},
	}

	for _, test := range tests {
		e := Extension{test.lf, nil}
		for _, obj := range test.objs {
			res := obj.obj.(filter.Labeled)
			if e.Transform(obj.obj); !reflect.DeepEqual(res.GetLabels(), obj.ans) {
				t.Error("transform did not modify obj correctly: ", res.GetLabels(), obj.ans)
			}
		}
	}
}

func TestBadTransform(t *testing.T) { // should not crash
	Extension{map[string]map[string]string{"x": map[string]string{"x": "x"}}, nil}.Transform(1)
}

func TestEmptyExtension(t *testing.T) {
	var lists = [][]string{
		[]string{"abc", "def", "ghi", "jkl"},
		[]string{""},
		[]string{"wow"},
		[]string{},
	}

	equal := func(a []string, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		sort.Strings(a)
		sort.Strings(b)
		return reflect.DeepEqual(a, b)
	}

	for _, l := range lists {
		e := EmptyExtension(l)
		l1 := &filter.LabelMatchEq{"any", "bad"}
		if !equal(l, e.FilterContexts(l, l1)) {
			t.Error("empty extension failed FilterContexts with LabelMatchEq on list", l)
		}
		if !equal(l, e.GetFilteredContexts(l1)) {
			t.Error("empty extension failed GetFilteredContexts with LabelMatchEq on list", l)
		}
		l2 := &filter.LabelMatchNeq{"any", "bad"}
		if !equal(l, e.FilterContexts(l, l2)) {
			t.Error("empty extension failed FilterContexts with LabelMatchNeq on list", l)
		}
		if !equal(l, e.GetFilteredContexts(l2)) {
			t.Error("empty extension failed GetFilteredContexts with LabelMatchNeq on list", l)
		}
	}
}
