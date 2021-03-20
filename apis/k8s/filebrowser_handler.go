package k8s

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubefilebrowser/apis"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/copyer"
	"kubefilebrowser/utils/execer"
	"kubefilebrowser/utils/logs"
	"strings"
)

type FileBrowserQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Pods      string `json:"pods" form:"pods" binding:"required"`
	Container string `json:"container" form:"container" binding:"required"`
	Path      string `json:"path" form:"path" binding:"required"`
}

// @Summary FileBrowser
// @description 容器文件浏览器
// @Tags Kubernetes
// @Param namespace query FileBrowserQuery true "namespace"
// @Param pods query FileBrowserQuery true "Pod名称"
// @Param container query FileBrowserQuery true "容器名称"
// @Param path query FileBrowserQuery true "路径"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/file_browser [get]
func FileBrowser(c *gin.Context) {
	render := apis.Gin{C: c}
	// {"ls", "-lQ", "--color=never", "--full-time", "/"}
	var query = &FileBrowserQuery{
		Path: "/",
	}
	if err := c.ShouldBindQuery(query); err != nil {
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	// check namespace
	_, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	pod, err := configs.RestClient.CoreV1().Pods(query.Namespace).
		Get(context.TODO(), query.Pods, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	var isUnix = true
	for _, value := range pod.Spec.NodeSelector {
		if strings.Contains(value, "windows") {
			isUnix = false
			break
		}
	}
	if isUnix {
		_, err = query.exec([]string{"sh"})
	} else {
		_, err = query.exec([]string{"cmd"})
	}
	if err != nil {
		logs.ErrorWithFields(err, logs.Fields{
			"annotation": "test container terminal shell",
		})
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	lsPath := "/ls_linux_amd64"
	command := []string{"/ls", query.Path}
	if !isUnix {
		lsPath = "/ls_windows_amd64.exe"
	}
	resByte, err := query.exec(command)
	if err != nil {
		logs.Error(err)
		if strings.Contains(err.Error(), "ls") ||
			err.Error() == "command terminated with exit code 126" {
			err = query.copyLsTar(lsPath)
			if err != nil {
				logs.Error(err)
				render.SetError(utils.CODE_ERR_APP, err)
				return
			}
			if isUnix {
				_cmd := []string{"chmod", "+x", "/ls"}
				_, err = query.exec(_cmd)
				if err != nil {
					logs.Error(err)
					render.SetError(utils.CODE_ERR_APP, err)
					return
				}
			}
			resByte, err = query.exec(command)
			if err != nil {
				logs.Error(err)
				render.SetError(utils.CODE_ERR_APP, err)
				return
			}
		} else {
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
	}
	var res []utils.File
	if err := json.Unmarshal(resByte, &res); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetRes(res, nil, 0)
}

func (query *FileBrowserQuery) copyLs(lsBinary, lsPath string) error {
	reader, writer := io.Pipe()
	cp := copyer.NewCopyer(query.Namespace, query.Pods, query.Container, configs.KuBeResConf, configs.RestClient)
	cp.Stdin = reader

	go func() {
		defer writer.Close()
		err := utils.MakeTar(lsBinary, lsPath, writer)
		if err != nil {
			logs.Error(err)
			return
		}
	}()

	return cp.CopyToPod(lsPath)
}

func (query *FileBrowserQuery) exec(command []string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	exec := execer.NewExec(query.Namespace, query.Pods, query.Container, configs.KuBeResConf, configs.RestClient)
	exec.Command = command
	exec.Tty = false
	exec.Stdout = &stdout
	exec.Stderr = &stderr
	err := exec.Exec()
	if err != nil {
		return nil, err
	}
	if stderr.String() != "" {
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
		err := utils.TarLs(lsPath, writer)
		if err != nil {
			logs.Error(err)
		}
	}()
	return cp.CopyToPod(lsPath)
}
