package copyer

import (
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type copyer struct {
	K8sClient     *kubernetes.Clientset
	RESTConfig    *rest.Config
	Namespace     string
	PodName       string
	ContainerName string
	Command       []string
	NoPreserve    bool
	Stdin         io.Reader
	Stdout        io.Writer
}
