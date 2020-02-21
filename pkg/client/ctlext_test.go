package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"testing"
)

func generateCtlExt(data map[string]string) runtime.Object {
	temp := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "configmap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ctl-config",
			Namespace: "kube-system",
		},
		Data: data,
	}
	return temp.DeepCopyObject()
}

func TestGetCtlExt(t *testing.T) {
	c := GetFakeConfigClient(map[string][]runtime.Object{
		"cluster1": {generateCtlExt(map[string]string{"_hidden": "true"})},
		"cluster2": {},
		"c4":       {generateCtlExt(make(map[string]string))},
	})
	// Hacky
	c.contextsGetter = StaticContextsGetter{contexts: []string{"cluster1", "cluster2", "cluster3", "c4"}}

	ans := map[string]map[string]string{
		"cluster1": {"_hidden": "true"},
		"cluster2": nil,
		"c4":       make(map[string]string),
	}

	if cf := c.GetCtlExt(); !reflect.DeepEqual(cf, ans) {
		t.Error("call to GetCtlExt expected", ans, "but instead returned", cf)
	}
}
