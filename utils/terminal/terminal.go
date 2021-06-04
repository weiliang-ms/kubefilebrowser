package terminal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
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
	msg       *utils.WsMessage
	xtermMsg  utils.XtermMessage
	copyData  []byte
	charEnter = []byte("\r")
)

// ssh流式处理器
type StreamHandler struct {
	SessionID   string
	WsConn      *utils.WsConnection
	ResizeEvent chan remotecommand.TerminalSize
	Shell       string
	inputCh     chan []byte
	outputCh    chan []byte
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

// Write executor回调向web端输出
func (handler *StreamHandler) Write(p []byte) (size int, err error) {
	// 产生副本
	copyData = make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	handler.outputCh <- copyData
	err = handler.WsConn.WsWrite(websocket.TextMessage, []byte(base64.StdEncoding.EncodeToString(copyData)))
	return
}

// Next executor回调获取web是否resize
func (handler *StreamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.ResizeEvent
	size = &ret
	return
}

// Read executor回调读取web端的输入
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
		handler.inputCh <- []byte(xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

func (handler *StreamHandler) Close() error {
	close(handler.outputCh)
	close(handler.inputCh)
	return nil
}

type CmdStruct struct {
	CommandID string
	Mode      string
	Msg       []byte
	PS1       string
}

func (handler *StreamHandler) CommandRecordChan() {
	var cmdlineCh = make(chan CmdStruct)
	handler.inputCh = make(chan []byte, 10)
	handler.outputCh = make(chan []byte, 10)
	go handler.splitCmdStream(cmdlineCh)
	for cmdLine := range cmdlineCh {
		switch cmdLine.Mode {
		case "input":
			input := ParseCmdInput(cmdLine.Msg)
			//input = strings.TrimPrefix(input, cmdLine.PS1)
			if len(input) == 0 {
				continue
			}
			logs.InfoWithFields(input, logs.Fields{
				"type": cmdLine.Mode,
				"uuid": cmdLine.CommandID,
			})
		case "output":
			output := ParseCmdOutput(cmdLine.Msg)
			//output = strings.TrimSuffix(output, cmdLine.PS1)
			if len(output) == 0 {
				continue
			}
			logs.InfoWithFields(output, logs.Fields{
				"type": cmdLine.Mode,
				"uuid": cmdLine.CommandID,
			})
		}
	}
}

func (handler *StreamHandler) splitCmdStream(cmdlineCh chan CmdStruct) {
	var (
		buf        bytes.Buffer
		ps1        string
		first      = true
		isEnter    bool
		inputState bool
		commandID  = uuid.NewV4().String()
	)
	// 回放功能
	replayRecorder := NewReplyRecord(handler.SessionID, handler.Shell)
	for {
		select {
		case <-handler.WsConn.CloseChan:
			replayRecorder.End()
			return
		case input := <-handler.inputCh:
			// 用户按下回车
			inputState = true
			isEnter = bytes.Contains(input, charEnter)
		case output := <-handler.outputCh:
			replayRecorder.Record(output)
			// 暂时不支持win系统的命令行收集
			if len(output) == 0 || handler.Shell == "cmd" || handler.Shell == "powershell" {
				continue
			}
			buf.Write(output)
			// 用户输入了Enter，开始结算命令
			if isEnter {
				// 处理异常情况， 上一次输出粘连本次输入
				if bytes.Contains(buf.Bytes(), []byte(ps1)) {
					_tmp := bytes.Split(buf.Bytes(), []byte(ps1))
					buf.Reset()
					buf.Write(bytes.Join(_tmp[1:], []byte("")))
				}
				commandID = uuid.NewV4().String()
				cmdlineCh <- CmdStruct{
					CommandID: commandID,
					Mode:      "input",
					PS1:       ps1,
					Msg:       buf.Bytes(),
				}
				isEnter, inputState = false, false
				buf.Reset()
			} else {
				// 用户又开始输入，并上次不处于输入状态，开始结算上次命令的结果
				if inputState {
					continue
				}
				// 只处理第一次
				if first {
					if s := GetPs1(output); s != "" {
						ps1 = s
					}
					first = false
					buf.Reset()
					continue
				}
				// 结算命令执行的结果
				if bytes.Contains(buf.Bytes(), []byte(ps1)) {
					if s := GetPs1(output); s != "" {
						ps1 = s
					}
					fmt.Println(ps1)
					cmdlineCh <- CmdStruct{
						CommandID: commandID,
						Mode:      "output",
						PS1:       ps1,
						Msg:       buf.Bytes(),
					}
					// 清空
					buf.Reset()
				}
			}
		}
	}
}
