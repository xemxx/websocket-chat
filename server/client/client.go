package client

import (
	"encoding/json"
	// "github.com/satori/go.uuid"

	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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
	uuid string
	conn *websocket.Conn
	send chan []byte

	stopCh chan struct{}
}

type PullMsg struct {
	Type    string `json:"type"`
	Uuid    string `json:"uid"`
	Message string `json:"message"`
	Bg      string `json:"bg,omitempty"`
}

type PushMsg struct {
	Err     bool   //是否错误
	Code    int    //错误代码
	Uuid    string `json:"uid"`
	Message string `json:"message,omitempty"`
	Bg      string `json:"bg,omitempty"`
}

//TODO: 声明消息json格式

//建立websockt连接
func HandleWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		zap.L().Error("upgrader error", zap.Error(err))
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256), stopCh: make(chan struct{})}
	manager.register <- client

	go client.handleMsg()
	go client.pushMsg()

}

func (c *Client) Stop() {
	close(c.stopCh)
	c.conn.Close()
}

//推送消息
func (c *Client) pushMsg() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case <-c.stopCh:
			return
		case msg, ok := <-c.send: //等待信号道
			if !ok {
				// 聊天室主动关闭信号道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			//TODO: 对消息进行判定然后重写消息发送的json格式并发送
			_, err = w.Write(msg)
			if err != nil {
				zap.L().Error(err.Error())
			}
			w.Write(newline)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

//拉取消息
func (c *Client) handleMsg() {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msgJson, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			return
		}
		msg := PullMsg{}

		if json.Unmarshal(msgJson, &msg) != nil {
			newSend := PushMsg{
				Err:     true,
				Code:    401,
				Message: "json error",
			}
			send, _ := json.Marshal(newSend)
			c.send <- send
			continue
		}
		zap.L().Debug("check type", zap.Any("msg", msg))
		switch msg.Type {
		case "bind":
			for client := range manager.clients {
				//删除其他已在线的连接
				if msg.Uuid == client.uuid {
					client.Stop()
					manager.unregister <- client
				}
			}
			c.uuid = msg.Uuid
			//反馈订阅成功
			newSend := PushMsg{
				Err:     false,
				Code:    200,
				Message: "bind success",
			}
			send, _ := json.Marshal(newSend)
			c.send <- send
		case "send":
			for client := range manager.clients {
				newSend := PushMsg{
					Err:     false,
					Code:    200,
					Message: msg.Message,
					Uuid:    msg.Uuid,
					Bg:      msg.Bg,
				}
				send, _ := json.Marshal(newSend)
				client.send <- send

			}
		}
	}
}

//TODO:实现日志记录
func checkMsgErr(err error) {
	if err != nil {
		fmt.Print(err)
		//TODO: 完善错误日志记录
	}
}
