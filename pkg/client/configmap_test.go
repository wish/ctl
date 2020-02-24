package client

import (
	"strconv"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getRandomConfigMaps(N int) []*corev1.ConfigMap {
	configmaps := make([]*corev1.ConfigMap, N)
	for n := 0; n < N; n++ {
		configmaps[n] = &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "configmap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
	}
	return configmaps
}

func getRandomConfigMapsObject(N int) []runtime.Object {
	configmaps := make([]runtime.Object, N)
	for n := 0; n < N; n++ {
		temp := &corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "configmap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
		configmaps[n] = temp.DeepCopyObject()
	}
	return configmaps
}

func TestListConfigMapsSingle(t *testing.T) {
	configmap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "configmap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": {configmap.DeepCopyObject()}})

	p, err := cl.ListConfigMaps("hi", "", ListOptions{})
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(p)

	if len(p) != 1 {
		t.Errorf("Unexpected number of configmaps found! %d != 1; expected", len(p))
	}
}

func TestListConfigMapsBadContext(t *testing.T) {
	configmap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "configmap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": {configmap.DeepCopyObject()}})

	p, err := cl.ListConfigMaps("pew", "", ListOptions{})

	if len(p) == 0 && err != nil {
		t.Log("Error as expected:", err.Error())
	} else {
		t.Error("Context not found did not error")
	}
}

func TestListConfigMapsMultiple(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6)})

	var queries = []struct {
		context   string
		namespace string
		size      int
	}{
		{"c1", "", 5},
		{"c2", "", 6},
		{"c1", "0", 1},
		{"c2", "1", 1},
	}

	for _, q := range queries {
		p, err := cl.ListConfigMaps(q.context, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of configmaps found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestListConfigMapsOverContexts(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6), "c3": getRandomConfigMapsObject(3), "c4": nil})

	var queries = []struct {
		contexts  []string
		namespace string
		size      int
	}{
		{[]string{"c1"}, "", 5},
		{[]string{"c2"}, "", 6},
		{[]string{"c3"}, "", 3},
		{nil, "", 14},
		{[]string{"c1", "c2", "c3"}, "", 14},
		{[]string{"c1", "c2"}, "5", 1},
		{nil, "4", 2},
		{[]string{}, "0", 3},
		{[]string{"c1", "c2"}, "", 11},
		{[]string{"c3", "c1"}, "", 8},
	}

	for _, q := range queries {
		p, err := cl.ListConfigMapsOverContexts(q.contexts, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of configmaps found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestGetConfigMap(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6)}
	cl := GetFakeConfigClient(vals)

	var queries = []struct {
		context   string
		namespace string
		name      string
	}{
		{"c1", "1", "1"},
		{"c2", "0", "0"},
		{"c1", "0", "0"},
		{"c2", "5", "5"},
	}

	for _, q := range queries {
		p, err := cl.GetConfigMap(q.context, q.namespace, q.name, GetOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if p.Name != q.name {
			t.Error("Incorrect configmap found")
		}
	}
}

func TestGetConfigMapBadContext(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6)}
	cl := GetFakeConfigClient(vals)

	var queries = []struct {
		context   string
		namespace string
		name      string
	}{
		{"c3", "", "1"},
		{"c3", "0", "0"},
		{"c3", "", "0"},
		{"c3", "5", "5"},
	}

	for _, q := range queries {
		c, err := cl.GetConfigMap(q.context, q.namespace, q.name, GetOptions{})
		if err != nil && c == nil {
			t.Log("Errored as expected:", err)
		} else {
			t.Error("Cronjob get did not error!")
		}
	}
}

func TestFindConfigMaps(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6), "c3": getRandomConfigMapsObject(3), "c4": nil})

	var queries = []struct {
		contexts  []string
		namespace string
		names     []string
		size      int
	}{
		{[]string{"c1"}, "", []string{"0"}, 1},
		{nil, "", []string{"0", "3", "5"}, 6},
		{nil, "4", []string{"0", "3", "5"}, 0},
		{nil, "4", []string{"4"}, 2},
		{[]string{"c2"}, "4", []string{"4"}, 1},
		{[]string{"c1", "c3", "c2"}, "4", []string{"4"}, 2},
		{[]string{"c3"}, "", []string{"3"}, 0},
	}

	for _, q := range queries {
		p, err := cl.FindConfigMaps(q.contexts, q.namespace, q.names, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of configmaps found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestFindConfigMapsError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6), "c3": getRandomConfigMapsObject(3)})

	var queries = []struct {
		contexts  []string
		namespace string
		names     []string
	}{
		{[]string{"c4"}, "", []string{"0"}},
		{[]string{"c4"}, "", []string{"0", "3", "5"}},
		{[]string{"c4"}, "4", []string{"0", "3", "5"}},
		{[]string{"c4"}, "4", []string{"4"}},
		{[]string{"c2", "c4"}, "4", []string{"4"}},
		{[]string{"c4", "c1", "c3", "c2"}, "4", []string{"4"}},
		{[]string{"c3", "c4"}, "", []string{"3"}},
	}

	for _, q := range queries {
		_, err := cl.FindConfigMaps(q.contexts, q.namespace, q.names, ListOptions{})
		if err == nil {
			t.Error("FindConfigMaps did not fail!")
		} else {
			t.Log("FindConfigMaps failed as expected: ", err.Error())
		}
	}
}

func TestFindConfigMap(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6), "c3": getRandomConfigMapsObject(3)})
	var queries = []struct {
		contexts  []string
		namespace string
		name      string
	}{
		{nil, "", "1"},
		{nil, "5", "5"},
		{[]string{"c2", "c3"}, "", "2"},
		{[]string{"c2", "c3"}, "3", "3"},
		{[]string{}, "", "2"},
	}

	for _, q := range queries {
		c, err := cl.findConfigMap(q.contexts, q.namespace, q.name, ListOptions{})

		if c == nil || err != nil {
			t.Error("Could not find pod with error:", err)
		} else if c.Name != q.name {
			t.Errorf("The found pod does not match the name requested: %s != %s", c.Name, q.name)
		}
	}
}

func TestFindConfigMapError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomConfigMapsObject(5), "c2": getRandomConfigMapsObject(6), "c3": getRandomConfigMapsObject(3)})
	var queries = []struct {
		contexts  []string
		namespace string
		name      string
	}{
		{[]string{"c2", "c3"}, "", "10"},
		{[]string{"c4"}, "3", "3"},
	}

	for _, q := range queries {
		_, err := cl.findConfigMap(q.contexts, q.namespace, q.name, ListOptions{})

		if err != nil {
			t.Log("Errored as expected:", err)
		} else {
			t.Error("Function did not error when finding log")
		}
	}
}
