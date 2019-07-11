package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strconv"
	"testing"
)

func getRandomPods(N int) []*corev1.Pod {
	pods := make([]*corev1.Pod, N)
	for n := 0; n < N; n++ {
		pods[n] = &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
	}
	return pods
}

func getRandomPodsObject(N int) []runtime.Object {
	pods := make([]runtime.Object, N)
	for n := 0; n < N; n++ {
		temp := &corev1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
		pods[n] = temp.DeepCopyObject()
	}
	return pods
}

func TestListPodsSingle(t *testing.T) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": []runtime.Object{pod.DeepCopyObject()}})

	p, err := cl.ListPods("hi", "", ListOptions{})
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(p)

	if len(p) != 1 {
		t.Errorf("Unexpected number of pods found! %d != 1; expected", len(p))
	}
}

func TestListPodsBadContext(t *testing.T) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": []runtime.Object{pod.DeepCopyObject()}})

	p, err := cl.ListPods("pew", "", ListOptions{})

	if len(p) == 0 && err != nil {
		t.Log("Error as expected:", err.Error())
	} else {
		t.Error("Context not found did not error")
	}
}

func TestListPodsMultiple(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomPodsObject(5), "c2": getRandomPodsObject(6)})

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
		p, err := cl.ListPods(q.context, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of pods found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestListPodsOverContexts(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomPodsObject(5), "c2": getRandomPodsObject(6), "c3": getRandomPodsObject(3), "c4": nil})

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
		p, err := cl.ListPodsOverContexts(q.contexts, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of pods found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestGetPod(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomPodsObject(5), "c2": getRandomPodsObject(6)}
	cl := GetFakeConfigClient(vals)

	var queries = []struct {
		context   string
		namespace string
		name      string
	}{
		{"c1", "", "1"},
		{"c2", "0", "0"},
		{"c1", "", "0"},
		{"c2", "5", "5"},
	}

	for _, q := range queries {
		p, err := cl.GetPod(q.context, q.namespace, q.name, GetOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if p.Name != q.name {
			t.Error("Incorrect pod found")
		}
	}
}

func TestFindPods(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomPodsObject(5), "c2": getRandomPodsObject(6), "c3": getRandomPodsObject(3), "c4": nil})

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
		p, err := cl.FindPods(q.contexts, q.namespace, q.names, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of pods found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestFindPodsError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomPodsObject(5), "c2": getRandomPodsObject(6), "c3": getRandomPodsObject(3)})

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
		_, err := cl.FindPods(q.contexts, q.namespace, q.names, ListOptions{})
		if err == nil {
			t.Error("FindPods did not fail!")
		} else {
			t.Log("FindPods failed as expected: ", err.Error())
		}
	}
}
