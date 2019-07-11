package cmd

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestDescribeSingle(t *testing.T) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := client.GetFakeConfigClient(map[string][]runtime.Object{"hi": []runtime.Object{pod.DeepCopyObject()}})

	cmd := GetDescribeCmd(cl)
	cmd.Flags().StringSliceP("context", "x", nil, "Context")
	cmd.SetArgs([]string{"test"})

	_, err := cmd.ExecuteC()

	fmt.Println(err)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDescribeBadContext(t *testing.T) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := client.GetFakeConfigClient(map[string][]runtime.Object{"hi": []runtime.Object{pod.DeepCopyObject()}})

	cmd := GetDescribeCmd(cl)
	cmd.Flags().StringSliceP("context", "x", nil, "Context")
	cmd.SetArgs([]string{"test", "--context=wow"})

	_, err := cmd.ExecuteC()

	fmt.Println(err)
	if err == nil {
		t.Error("Was expecting cmd execution to error")
	} else {
		t.Log("Error as expected:", err)
	}
}

func TestDescribeUnfound(t *testing.T) {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := client.GetFakeConfigClient(map[string][]runtime.Object{"hi": []runtime.Object{pod.DeepCopyObject()}})

	cmd := GetDescribeCmd(cl)
	cmd.Flags().StringSliceP("context", "x", nil, "Context")
	cmd.SetArgs([]string{"pew"})

	_, err := cmd.ExecuteC()

	fmt.Println(err)
	if err == nil {
		t.Error("Was expecting cmd execution to error")
	} else {
		t.Log("Error as expected:", err)
	}
}
