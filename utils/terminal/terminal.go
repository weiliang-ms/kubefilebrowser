package kube

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
)

// Terminal operating start
var (
	msg      *utils.WsMessage
	xtermMsg utils.XtermMessage
	copyData []byte
)

// ssh流式处理器
type StreamHandler struct {
	WsConn      *utils.WsConnection
	ResizeEvent chan remotecommand.TerminalSize
}

type WebTerminal struct {
	K8sClient *kubernetes.Clientset
	PodName   string
	Namespace string
	Container string
	Shell     string
}

func (wt *WebTerminal) NewSshSPDYExecutor() *rest.Request {
	return wt.K8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(wt.PodName).
		Namespace(wt.Namespace).
		SubResource("exec").
		VersionedParams(&coreV1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Command:   []string{wt.Shell},
			Container: wt.Container,
		}, scheme.ParameterCodec)
}

// executor回调向web端输出
func (handler *StreamHandler) Write(p []byte) (size int, err error) {

	// 产生副本
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.WsConn.WsWrite(websocket.TextMessage, []byte(base64.StdEncoding.EncodeToString(copyData)))
	return
}

// executor回调获取web是否resize
func (handler *StreamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.ResizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *StreamHandler) Read(p []byte) (size int, err error) {
	// 读web发来的输入
	if msg, err = handler.WsConn.WsRead(); err != nil {
		return
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(string(msg.Data))
	if err != nil {
		logs.Error("websock cmd string base64 decoding failed", err)
		return 0, err
	}
	if err = json.Unmarshal(decodeBytes, &xtermMsg); err != nil {
		return
	}
	//web ssh调整了终端大小
	switch xtermMsg.Type {
	case utils.WsMsgResize:
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.ResizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	case utils.WsMsgInput: // web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return

}
