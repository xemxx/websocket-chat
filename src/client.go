package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"time"
	"bytes"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
    WriteBufferSize:  1024,
    HandshakeTimeout: 5 * time.Second,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func HandleWs(hub *Hub,w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.pushMsg()
	go client.pullMsg()
}

func (c *Client) pushMsg(){
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
			case msg,ok:=<-c.send://等待信号道
				if !ok{
					// 聊天室主动关闭信号道
					c.conn.WriteMessage(websocket.CloseMessage,[]byte{})
					return
				}
				w,err:=c.conn.NextWriter(websocket.TextMessage)
				if err != nil{
					return
				}
				w.Write(msg)

				//消息队列
				n := len(c.send)
				for i := 0; i < n; i++ {
					w.Write(newline)
					w.Write(<-c.send)
				}
				if err:=w.Close(); err!=nil{
					return
				}
		}
	}
}

func (c *Client) pullMsg(){
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_,msg,err:=c.conn.ReadMessage()
		if err!= nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg =bytes.TrimSpace(bytes.Replace(msg,newline,space,-1))
		c.hub.broadcast <- msg
	}
}