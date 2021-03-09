package execer

import (
	"io"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type execer struct {
	K8sClient     *kubernetes.Clientset
	RESTConfig    *rest.Config
	Namespace     string
	PodName       string
	ContainerName string
	Command       []string
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
	Tty           bool
}

func NewExec(namespace, podName, containerName string, restConfig *rest.Config, k8sClient *kubernetes.Clientset) *execer {
	return &execer{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
		RESTConfig:    restConfig,
		K8sClient:     k8sClient,
	}
}

// Exec 在给定容器中执行命令
func (e *execer) Exec() error {
	req := e.K8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(e.PodName).
		Namespace(e.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := coreV1.AddToScheme(scheme); err != nil {
		return err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&coreV1.PodExecOptions{
		Command:   e.Command,
		Container: e.ContainerName,
		Stdin:     e.Stdin != nil,
		Stdout:    e.Stdout != nil,
		Stderr:    e.Stderr != nil,
		TTY:       e.Tty,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(e.RESTConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	var sizeQueue remotecommand.TerminalSizeQueue
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:             e.Stdin,
		Stdout:            e.Stdout,
		Stderr:            e.Stderr,
		Tty:               e.Tty,
		TerminalSizeQueue: sizeQueue,
	})
	if err != nil {
		return err
	}
	return nil
}