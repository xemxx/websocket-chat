package client

import (
	"encoding/json"
	// "github.com/satori/go.uuid"
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
	uuid string
	conn *websocket.Conn
	send chan []byte
}

type PullMsg struct {
	Type     string `json:"type"`
	Uuid     string `json:"uid"`
	ToUuid	 string `json:"touid"`
	Message  string `json:"message"`
}

type PushMsg struct{
	Err  bool		//是否错误
	Code int		//错误代码
	Message string  //具体数据
}

//TODO: 声明消息json格式

//建立websockt连接
func HandleWs(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// u,_:=uuid.NewV4()
	// client := &Client{uuid: u.Bytes(), conn: conn, send: make(chan []byte, 256)}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	manager.register <- client

	//TODO: 初始读取对应用户是否有未读消息，并循环推送消息
	go client.pushMsg()
	go client.pullMsg()
	
	// //TODO:主动推送uuid让前台绑定用户返回用户id
	// client.send<-client.uuid
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
				//TODO: 对消息进行判定然后重写消息发送的json格式并发送
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
		_,msgJson,err:=c.conn.ReadMessage()
		if err!= nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}
		//TODO: 解析msg的json finish
		msg:= PullMsg{}
		if json.Unmarshal(msgJson,&msg) != nil{
			//TODO:返回错误信息
			
			return
		}

		//TODO: 通过账户寻找发送账户是否在线并推送  finish
		switch msg.Type{
			case "bind":
				c.uuid=msg.Uuid
				//TODO:查询数据库是否有未读消息如有则推送
			case "send":
				for client:=range manager.clients{
					if msg.ToUuid==client.uuid{
						//TODO:按照固定json传输 finish
						newSend:=PushMsg{
							Err:false,
							Code:200,
							Message:msg.Message,
						}
						send,err:=json.Marshal(newSend)
						if err !=nil{
							//TODO:记录错误
						}
						client.send<-send
					}
				}
		}

		//TODO: 保存读取到的消息到数据库做聊天记录


		//msg =bytes.TrimSpace(bytes.Replace(msg,newline,space,-1))
		//c.manager.broadcast <- msg
	}
}