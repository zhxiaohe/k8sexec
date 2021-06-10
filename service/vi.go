package service

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsMessage struct {
	MessageType int
	Data        []byte
}

type WsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *WsMessage // 读取队列
	outChan  chan *WsMessage // 发送队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

func (ws *WsConnection) readLoop() {
	for {
		msgType, data, err := ws.wsSocket.ReadMessage()
		if err != nil {
			ws.wsClose()
			return
		}
		msg := &WsMessage{
			msgType,
			data,
		}
		// 读取websocket客户端消息，写入chan
		select {
		case ws.inChan <- msg:
		case <-ws.closeChan:
			fmt.Println("recvwebsocket3:")
			ws.wsClose()
		}
	}
}

func (ws *WsConnection) writeLoop() {
	for {
		select {
		// 读取outChan队列，返回至websockt客户端
		case msg := <-ws.outChan:
			err := ws.wsSocket.WriteMessage(msg.MessageType, msg.Data)
			if err != nil {
				ws.wsClose()
			}
		case <-ws.closeChan:
			ws.wsClose()
		}
	}

}

func (ws *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	//获取websocket inchan消息，return to k8sexec终端
	case msg = <-ws.inChan:
		fmt.Println("read: ", msg)
		return msg, nil
	case <-ws.closeChan:
		return nil, errors.New("websocket closed")
	}

}

func (ws *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	//获取k8sexec终端消息至队列
	case ws.outChan <- &WsMessage{messageType, data}:
		return nil
	case <-ws.closeChan:
		return errors.New("websocket closed")
	}
}

func (ws *WsConnection) wsClose() {
	ws.wsSocket.Close()
}

func wsInit(c *gin.Context) (*WsConnection, error) {
	var ws *websocket.Conn
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
	}
	WS := &WsConnection{
		wsSocket:  ws,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	go WS.readLoop()
	go WS.writeLoop()
	return WS, nil

}

func Echo(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
