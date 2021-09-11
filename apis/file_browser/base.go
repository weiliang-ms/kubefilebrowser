package file_browser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/copyer"
	"kubefilebrowser/utils/execer"
	"kubefilebrowser/utils/logs"
	"strings"
)

type FileBrowserQuery struct {
	Namespace string    `json:"namespace" form:"namespace" binding:"required"`
	Pods      string    `json:"pods" form:"pods" binding:"required" binding:"required"`
	Container string    `json:"container" form:"container" binding:"required"`
	Path      string    `json:"path" form:"path" binding:"required" binding:"required"`
	OldPath   string    `json:"old_path,omitempty" form:"old_path"`
	Command   []string  `json:"-"`
	Stdin     io.Reader `json:"-"`
}

func (query *FileBrowserQuery) FileBrowser() (res []byte, err error) {
	// check namespace
	_, err = configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	// check pod
	pod, err := configs.RestClient.CoreV1().Pods(query.Namespace).
		Get(context.TODO(), query.Pods, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	var osType = "linux"
	var arch = "amd64"

	// get pod system arch and type
	node, err := configs.RestClient.CoreV1().Nodes().
		Get(context.TODO(), pod.Spec.NodeName, metaV1.GetOptions{})
	if err == nil {
		if node.Labels["beta.kubernetes.io/os"] != "" {
			osType = node.Labels["beta.kubernetes.io/os"]
		} else if node.Labels["kubernetes.io/os"] != "" {
			osType = node.Labels["kubernetes.io/os"]
		}
		if node.Labels["beta.kubernetes.io/arch"] != "" {
			arch = node.Labels["beta.kubernetes.io/arch"]
		} else if node.Labels["kubernetes.io/arch"] != "" {
			arch = node.Labels["kubernetes.io/arch"]
		}
	}

	lsPath := fmt.Sprintf("/kf_tools_%s_%s", osType, arch)
	if osType == "windows" {
		lsPath = fmt.Sprintf("/kf_tools_%s_%s.exe", osType, arch)
	}
	reTryCmd := query.Command
	res, err = query.exec()
	if err != nil {
		logs.Error(err)
		if strings.Contains(err.Error(), "kf_tools") ||
			err.Error() == "command terminated with exit code 126" {
			if osType != "windows" {
				query.Command = []string{"sh"}
				_, err = query.exec()
			} else {
				query.Command = []string{"cmd"}
				_, err = query.exec()
			}
			if err != nil {
				logs.ErrorWithFields(err, logs.Fields{
					"annotation": "test container terminal shell",
				})
				logs.Error(err)
				return nil, err
			}
			err = query.copyLsTar(lsPath)
			if err != nil {
				logs.Error(err)
				return nil, err
			}
			if osType != "windows" {
				_cmd := []string{"chmod", "+x", "/kf_tools"}
				query.Command = _cmd
				_, err = query.exec()
				if err != nil {
					logs.Error(err)
					return nil, err
				}
			}
			query.Command = reTryCmd
			res, err = query.exec()
			if err != nil {
				logs.Error(err)
				logs.Error(err)
				return nil, err
			}
		} else {
			logs.Error(err)
			return nil, err
		}
	}
	return res, err
}

func (query *FileBrowserQuery) exec() ([]byte, error) {
	var stdout, stderr bytes.Buffer
	exec := execer.NewExec(query.Namespace, query.Pods, query.Container, configs.KuBeResConf, configs.RestClient)
	exec.Command = query.Command
	exec.Tty = false
	exec.Stdin = query.Stdin
	exec.Stdout = &stdout
	exec.Stderr = &stderr
	err := exec.Exec()
	if err != nil {
		if len(stderr.String()) != 0 {
			logs.Error(err)
			return nil, fmt.Errorf(stderr.String())
		} else {
			return nil, err
		}
	}
	if len(stderr.Bytes()) != 0 {
		return nil, fmt.Errorf(stderr.String())
	}
	return stdout.Bytes(), nil
}

func (query *FileBrowserQuery) copyLsTar(lsPath string) error {
	reader, writer := io.Pipe()
	cp := copyer.NewCopyer(query.Namespace, query.Pods, query.Container, configs.KuBeResConf, configs.RestClient)
	cp.Stdin = reader

	go func() {
		defer writer.Close()
		err := utils.TarKFTools(lsPath, writer)
		if err != nil {
			logs.Error(err)
		}
	}()
	return cp.CopyToPod(lsPath)
}
