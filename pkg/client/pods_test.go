package client

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestListPods(t *testing.T) {
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
