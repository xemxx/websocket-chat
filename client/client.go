package client

import (
	"encoding/json"
	// "github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"time"
	"websocket-chat/mysql"
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
	join string
}

type PullMsg struct {
	Type     string `json:"type"`
	Uuid     string `json:"uid"`
	ToUuid	 string `json:"touid"`
	Message  string `json:"msg"`
}

type PushMsg struct{
	Type 	string 	`json:"type"`
	Err  	bool	`json:"error"`
	Code 	int		`json:"code"`
	Message string  `json:"msg"`
	Uuid 	string 	`json:"uid"`
	ToUuid 	string	`json:"touid"`
}

//TODO: 声明消息json格式

//建立websockt连接
func HandleWs(w http.ResponseWriter, r *http.Request){
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	manager.register <- client

	go client.pushMsg()
	go client.pullMsg()
	
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

//拉取消息
func (c *Client) pullMsg(){
	db := mysql.NewMysql()
	defer func() {
		manager.unregister <- c
		c.conn.Close()
		db.Close()
	}()
	for {
		_,msgJson,err:=c.conn.ReadMessage()
		if err!= nil{
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			return
		}
		//TODO: 解析msg的json finish
		msg:= PullMsg{}
		if json.Unmarshal(msgJson,&msg) != nil{
			//TODO:返回错误信息  finish 
			c.send<-sendMsg(&PushMsg{"",true,401,"json error","",""})
			continue
		}
		if (msg.Uuid!=c.uuid) && c.isBind(){
			c.send<-sendMsg(&PushMsg{"",true,405,"uid error","",""})
			continue
		}
		//TODO: 通过账户寻找发送账户是否在线并推送  finish
		switch msg.Type{
			case "bind":
				if c.isBind(){
					c.send<-sendMsg(&PushMsg{msg.Type,true,403,"请勿重复绑定","",""})
					continue
				}
				for client:=range manager.clients{
					//删除其他已在线的连接
					if msg.Uuid==client.uuid{
						client.conn.Close()
						manager.unregister <- client
						//TODO:改为redis的方式解决登录和注销问题
					}
				}
				c.uuid=msg.Uuid
				//反馈订阅成功
				c.send<-sendMsg(&PushMsg{msg.Type,false,201,"bind success","",""})

			case "send":
				if !c.isBind(){
					c.send<-sendMsg(&PushMsg{msg.Type,true,402,"请先绑定后请求","",""})
					continue
				}
				is_read:=0
				for client:=range manager.clients{
					//判断是否在线
					if msg.ToUuid==client.uuid{
						//TODO:按照固定json传输 finish
						client.send<-sendMsg(&PushMsg{msg.Type,false,200,msg.Message,msg.Uuid,""})
						if client.join == msg.ToUuid{
							is_read=1
						}
						break
					}
				}
				//TODO: 保存读取到的消息到数据库做聊天记录  finish
				stmt,err:=db.Prepare("insert into msg(uid,touid,send_time,is_read,msg)values(?,?,?,?,?)")
				if err != nil {
					fmt.Print(err)
					continue
				}
				_,err=stmt.Exec(c.uuid,msg.ToUuid,time.Now().Unix(),is_read,msg.Message)
				if err != nil {
					fmt.Print(err)
				}
				stmt.Close()
			case "join":
				if c.isBind(){
					c.send<-sendMsg(&PushMsg{msg.Type,true,402,"请先绑定后请求","",""})
					continue
				}
				c.join=msg.ToUuid
				c.send<-sendMsg(&PushMsg{msg.Type,false,202,"join success","",""})

				rows,err:=db.Query("update msg set (is_read=0) where is_read=0 and uid=? and touid=?",msg.ToUuid,msg.Uuid)
				if err != nil {
					fmt.Print(err)
					continue
				}
				rows.Close()
			case "exit":
				if c.isBind(){
					c.send<-sendMsg(&PushMsg{msg.Type,true,402,"请先绑定后请求","",""})
					continue
				}
				c.join=""
				c.send<-sendMsg(&PushMsg{msg.Type,false,200,"exit success","",""})
		}
	}
}

func (c *Client)isBind()bool{
	if c.uuid==""{
		return false
	}
	return true;
}

func sendMsg(newSend *PushMsg)[]byte{
	send,_:=json.Marshal(newSend)
	return send
}
