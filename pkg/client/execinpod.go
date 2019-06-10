package client

import (
  "io"
  corev1 "k8s.io/api/core/v1"
  "k8s.io/client-go/tools/remotecommand"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/apimachinery/pkg/runtime"
  "github.com/ContextLogic/ctl/pkg/client/helper"
)

// TODO: Add options param???
func (c *Client) ExecInPod(contexts []string, namespace, name, container string, commands []string, stdin io.Reader, stdout, stderr io.Writer) error {
  pod, container, err := c.findPodWithContainer(contexts, namespace, name, container)
  if err != nil {
    return err
  }

  cl, err := c.getContextClientset(pod.Context)
  if err != nil {
    panic(err.Error())
  }

  // Credit to a4abhishek for most of this https://github.com/a4abhishek/Client-Go-Examples/blob/master/exec_to_pod/exec_to_pod.go
  req := cl.CoreV1().RESTClient().Post().
    Resource("pods").
    Name(pod.Name).
    Namespace(pod.Namespace).
    SubResource("exec")

  scheme := runtime.NewScheme()
  if err := corev1.AddToScheme(scheme); err != nil {
    return err
  }

  req.VersionedParams(&corev1.PodExecOptions{
    Container: container,
    Command: commands, // COMMAND
    Stdin: stdin != nil,
    Stdout: true,
    Stderr: true,
    TTY: false,
  }, runtime.NewParameterCodec(scheme))

  config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
    &clientcmd.ClientConfigLoadingRules{ExplicitPath: helper.GetKubeConfigPath()},
    &clientcmd.ConfigOverrides{CurrentContext: pod.Context}).ClientConfig()

  exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
  if err != nil {
    return err
  }

  err = exec.Stream(remotecommand.StreamOptions{
    Stdin: stdin,
    Stdout: stdout,
    Stderr: stderr,
    Tty: false,
  })
  if err != nil {
    return err
  }
  return nil
}
