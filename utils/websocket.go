package utils

import (
	"encoding/base64"
	"errors"
	"github.com/gorilla/websocket"
	"kubefilebrowser/utils/logs"
	"net/http"
	"sync"
	"time"
)

// http升级websocket协议的配置
var wsUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket消息
type WsMessage struct {
	MessageType int
	Data        []byte
}

// 封装websocket连接
type WsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *WsMessage // 读取队列
	outChan  chan *WsMessage // 发送队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

// 读取协程
func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		// 读一个message
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERROR
		}
		msg = &WsMessage{
			msgType,
			data,
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- msg:
		case <-wsConn.closeChan:
			goto ERROR
		}
	}
ERROR:
	wsConn.WsClose()
}

// 发送协程
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		// 取一个应答
		case msg = <-wsConn.outChan:
			// 写给websocket
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto ERROR
			}
		case <-wsConn.closeChan:
			goto ERROR
		}
	}
ERROR:
	wsConn.WsClose()
}

// 检查客户端的存活
func (wsConn *WsConnection) procLoop() {
	// 启动一个gouroutine发送心跳
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if err := wsConn.WsWrite(websocket.PingMessage, []byte("heartbeat from server")); err != nil {
				goto ERROR
			}
		}
	ERROR:
		wsConn.WsClose()
	}()
}

/************** 并发安全 API **************/
func InitWebsocket(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	// 应答客户端告知升级连接为websocket
	if wsSocket, err = wsUpgrader.Upgrade(resp, req, nil); err != nil {
		return
	}
	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 1024),
		outChan:   make(chan *WsMessage, 1024),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	logs.Info(wsSocket.RemoteAddr(), " WebSocket connection")
	// 存活检测
	go wsConn.procLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()

	return
}

// 发送消息
func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:
	case <-wsConn.closeChan:
		err = errors.New("WebSocket connection closed")
	}
	return
}

// 读取消息
func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.closeChan:
		err = errors.New("WebSocket connection closed")
	}
	return
}

// 关闭连接
func (wsConn *WsConnection) WsClose() {
	_ = wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		logs.Warn(wsConn.wsSocket.RemoteAddr(), " WebSocket connection closed")
		//<-wsConn.closeChan
		close(wsConn.closeChan)
	}
	wsConn.mutex.Unlock()
	return
}

// 错误检查
func WsHandleError(ws *WsConnection, err error) bool {
	if err != nil {
		dt := time.Now().Add(time.Second)
		if err := ws.wsSocket.WriteMessage(websocket.TextMessage, []byte(base64.StdEncoding.EncodeToString([]byte(err.Error())))); err != nil {
			logs.Error(ws.wsSocket.RemoteAddr(), err)
		}
		if err := ws.wsSocket.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logs.Error(ws.wsSocket.RemoteAddr(), err)
		}
		return true
	}
	return false
}

// web终端发来的包
type XtermMessage struct {
	Type  string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input string `json:"input"` // msgtype=input情况下使用
	Rows  uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols  uint16 `json:"cols"`  // msgtype=resize情况下使用
}

const (
	WsMsgInput  = "input"
	WsMsgResize = "resize"
)
