package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
	"kubecp/configs"
	"kubecp/controller"
	"kubecp/logs"
	"kubecp/utils"
	kube "kubecp/utils/terminal"
)

type TerminalQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	Pods      string `json:"pods" form:"pods" binding:"required"`
	Container string `json:"container" form:"container" binding:"required"`
	Shell     string `json:"shell" form:"shell"`
}

// @Summary Container terminal
// @description pod 中容器的终端
// @Tags Kubernetes
// @Param namespace query TerminalQuery true "namespace" default(default)
// @Param pods query TerminalQuery true "Pod名称"
// @Param container query TerminalQuery true "容器名称"
// @Param shell query TerminalQuery false "shell" default(sh[bash/sh/cmd])
// @Success 200 {object} controller.JSONResult
// @Failure 500 {object} controller.JSONResult
// @Router /api/k8s/terminal [get]
func Terminal(c *gin.Context) {
	render := controller.Gin{C: c}
	var query = &TerminalQuery{
		Shell: "sh",
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
	_, err = configs.RestClient.CoreV1().Pods(query.Namespace).
		Get(context.TODO(), query.Pods, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	wsConn, err := utils.InitWebsocket(c.Writer, c.Request)
	if utils.WsHandleError(wsConn, err) {
		wsConn.WsClose()
		return
	}
	defer wsConn.WsClose()
	webTerminal := kube.WebTerminal{
		K8sClient: configs.RestClient,
		Namespace: query.Namespace,
		PodName:   query.Pods,
		Container: query.Container,
		Shell:     query.Shell,
	}
	SshSPDYExecutor := webTerminal.NewSshSPDYExecutor()
	executor, err := remotecommand.NewSPDYExecutor(configs.KuBeResConf, "POST", SshSPDYExecutor.URL())
	if utils.WsHandleError(wsConn, err) {
		wsConn.WsClose()
		return
	}
	handler := &kube.StreamHandler{WsConn: wsConn, ResizeEvent: make(chan remotecommand.TerminalSize)}
	if utils.WsHandleError(wsConn, executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		Tty:               true,
		TerminalSizeQueue: handler,
	})) {
		wsConn.WsClose()
	}
}
