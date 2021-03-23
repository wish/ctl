package client

import (
	"io"

	"github.com/wish/ctl/pkg/client/helper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecInPod executes a command on a pod interactively
func (c *Client) ExecInPod(contexts []string, namespace, name, container string, options ListOptions, commands []string, stdin io.Reader, stdout, stderr io.Writer) error {
	pod, container, err := c.FindPodWithContainer(contexts, namespace, name, container, options)
	if err != nil {
		return err
	}

	cl, err := c.getContextInterface(pod.Context)
	if err != nil {
		panic(err.Error())
	}

	// Credit to a4abhishek for most of this https://github.com/a4abhishek/Client-Go-Examples/blob/master/exec_to_pod/exec_to_pod.go
	req := cl.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	myscheme := runtime.NewScheme()
	if err := corev1.AddToScheme(myscheme); err != nil {
		return err
	}

	req.VersionedParams(&corev1.PodExecOptions{
		Container: container,
		Command:   commands, // COMMAND
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: helper.GetKubeConfigPath()},
		&clientcmd.ConfigOverrides{CurrentContext: pod.Context}).ClientConfig()

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    true,
	})
	if err != nil {
		return err
	}
	return nil
}
