package k8s

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"kubefilebrowser/apis"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/execer"
	"kubefilebrowser/utils/logs"
)

type ExecQuery struct {
	Namespace string   `json:"namespace" form:"namespace" binding:"required"`
	Pods      string   `json:"pods" form:"pods" binding:"required"`
	Container string   `json:"container" form:"container" binding:"required"`
	Stdout    bool     `json:"stdout" form:"stdout"`
	Stderr    bool     `json:"stderr" form:"stderr"`
	Tty       bool     `json:"tty" form:"tty"`
	Command   []string `json:"command" form:"command"`
}
// Exec
// @Summary Container 执行器
// @description 在pod的容器中执行
// @Tags Kubernetes
// @Param namespace query ExecQuery true "namespace"
// @Param pods query ExecQuery true "Pod名称"
// @Param container query ExecQuery true "容器名称"
// @Param command query ExecQuery true "命令"
// @Param stdout query ExecQuery false "标准输出"
// @Param stderr query ExecQuery false "错误输出"
// @Param tty query ExecQuery false "终端"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/exec [get]
func Exec(c *gin.Context) {
	render := apis.Gin{C: c}
	var query = &ExecQuery{
		Stdout:  true,
		Stderr:  true,
		Tty:     false,
		Command: []string{"ls", "-lQ", "--color=never", "--full-time", "/"},
	}
	if err := c.ShouldBindQuery(query); err != nil {
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	check := Check{
		namespace: query.Namespace,
		pod:       query.Pods,
	}
	// check namespace
	if _, err := check.Namespace(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	if _, err := check.Pod(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	exec := execer.NewExec(query.Namespace, query.Pods, query.Container, configs.KuBeResConf, configs.RestClient)
	exec.Command = query.Command
	exec.Tty = query.Tty
	var stdout, stderr bytes.Buffer
	if query.Stdout {
		exec.Stdout = &stdout
	}
	if query.Stderr {
		exec.Stderr = &stderr
	}
	if err := exec.Exec(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	render.SetJson(map[string]interface{}{
		"err": stderr.String(),
		"out": stdout.String(),
	})
}
