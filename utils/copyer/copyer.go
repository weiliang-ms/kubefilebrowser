package copyer

import (
	"bytes"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kubefilebrowser/utils/execer"
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
	var stderr bytes.Buffer
	exec := execer.NewExec(c.Namespace, c.PodName, c.ContainerName, c.RESTConfig, c.K8sClient)
	exec.Tty = false
	exec.Stderr = &stderr
	exec.Stdin = c.Stdin
	exec.Stdout = c.Stdout
	exec.Command = c.Command
	err := exec.Exec()
	if err != nil {
		return nil, err
	}
	return stderr.Bytes(), nil
}
