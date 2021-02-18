package copyer

import (
	"bytes"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func NewCopyer(namespace, podName, containerName string, restConfig *rest.Config, k8sClient *kubernetes.Clientset) *copyer {
	return &copyer{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		RESTConfig:    restConfig,
		K8sClient:     k8sClient,
	}
}

// Exec 在给定容器中执行命令
func (c *copyer) Exec() ([]byte, error) {
	req := c.K8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(c.PodName).
		Namespace(c.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := coreV1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&coreV1.PodExecOptions{
		Command:   c.Command,
		Container: c.ContainerName,
		Stdin:     c.Stdin != nil,
		Stdout:    c.Stdout != nil,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.RESTConfig, "POST", req.URL())
	if err != nil {
		return nil, err
	}

	var stderr bytes.Buffer
	var sizeQueue remotecommand.TerminalSizeQueue
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             c.Stdin,
		Stdout:            c.Stdout,
		Stderr:            &stderr,
		Tty:               false,
		TerminalSizeQueue: sizeQueue,
	})
	if err != nil {
		return nil, err
	}

	return stderr.Bytes(), nil
}
