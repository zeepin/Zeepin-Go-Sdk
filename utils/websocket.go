package utils

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

//WebSocketClient use for client to operation web socket
type WebSocketClient struct {
	addr      string
	conn      *websocket.Conn
	existCh   chan interface{}
	OnConnect func(address string)
	OnClose   func(address string)
	OnError   func(address string, err error)
	OnMessage func([]byte)
	lock      sync.RWMutex
	status    bool
}

//Create WebSocketClient instance
func NewWebSocketClient() *WebSocketClient {
	return &WebSocketClient{
		OnConnect: func(address string) {},
		OnClose:   func(address string) {},
		OnError:   func(address string, err error) {},
		OnMessage: func([]byte) {},
		existCh:   make(chan interface{}, 0),
	}
}

//Connect to server
func (this *WebSocketClient) Connect(addr string) (err error) {
	this.addr = addr
	this.conn, _, err = websocket.DefaultDialer.Dial(this.addr, nil)
	if err != nil {
		return err
	}
	if this.OnConnect != nil {
		this.OnConnect(this.addr)
	}
	this.status = true
	go this.doRecv()
	return nil
}

//Send data to server
func (this *WebSocketClient) Send(data []byte) error {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if !this.status {
		return fmt.Errorf("WebSocket connect has already closed.")
	}
	return this.conn.WriteMessage(websocket.TextMessage, data)
}

func (this *WebSocketClient) doRecv() {
	defer this.Close()
	for {
		_, data, err := this.conn.ReadMessage()
		if err != nil {
			if this.Status() && this.OnError != nil {
				this.OnError(this.addr, fmt.Errorf("ReadMessage error:%s", err))
			}
			return
		}
		if this.OnMessage != nil {
			this.OnMessage(data)
		}
	}
}

//Status return the status of connection of client and server
func (this *WebSocketClient) Status() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.status
}

//Close the connection of server
func (this *WebSocketClient) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if !this.status {
		return nil
	}
	this.status = false
	close(this.existCh)
	if this.OnClose != nil {
		this.OnClose(this.addr)
	}
	return this.conn.Close()
}
