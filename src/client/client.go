package client

import (
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"time"
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
	uuid  []byte
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

//TODO: 声明消息json格式

//建立websockt连接
func HandleWs(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	u,_:=uuid.NewV4()
	client := &Client{uuid: u.Bytes(), conn: conn, send: make(chan []byte, 256)}
	manager.register <- client

	//TODO: 初始读取对应用户是否有未读消息，并循环推送消息
	go client.pushMsg()
	go client.pullMsg()
	
	//TODO:主动推送uuid让前台绑定用户返回用户id
	client.send<-client.uuid
}

//推送消息
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
				//TODO: 重写消息发送的json格式并发送
				w.Write(msg)
				//消息队列，节约推送时间
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
		manager.unregister <- c
		c.conn.Close()
	}()
	for {
		//_,msg,err:=c.conn.ReadMessage()
		_,_,err:=c.conn.ReadMessage()
		if err!= nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}
		//TODO: 解析msg的json
		//TODO: 通过账户寻找发送账户是否在线并推送
		//TODO: 保存读取到的消息到数据库
		

		//msg =bytes.TrimSpace(bytes.Replace(msg,newline,space,-1))
		//c.manager.broadcast <- msg
	}
}