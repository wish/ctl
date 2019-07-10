package cmd

import (
	"github.com/ContextLogic/ctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestDescribeEmpty(t *testing.T) {
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

	cmd.Run(cmd, []string{"hi"})
}
