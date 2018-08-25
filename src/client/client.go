package client

import (
	"encoding/json"
	// "github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
	Uuid string 	//发送者
	ToUuid string	//接受者
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
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chat")
	if err != nil {
		fmt.Print(err)
		return 
	}
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
			newSend:=PushMsg{
				Err:true,
				Code:401,
				Message:"json error",
			}
			send,_:=json.Marshal(newSend)
			c.send<-send
			continue
		}

		//TODO: 通过账户寻找发送账户是否在线并推送  finish
		switch msg.Type{
			case "bind":
				c.uuid=msg.Uuid
				//TODO:查询数据库是否有未读消息如有则推送  finish
				rows,err:=db.Query("select uid,touid,msg from msg where is_read=? and uid=?",0,msg.Uuid)
				if err != nil {
					fmt.Print(err)
					continue
				}
				for rows.Next(){
					sendMsg:=new(PushMsg)
					err = rows.Scan(&sendMsg.Uuid, &sendMsg.ToUuid,&sendMsg.Message)
					sendMsg.Err=false
					sendMsg.Code=200
					send,_:=json.Marshal(*sendMsg)
					c.send<-send
				}
			case "send":
				is_read:=0
				for client:=range manager.clients{
					//判断是否在线
					if msg.ToUuid==client.uuid{
						//TODO:按照固定json传输 finish
						newSend:=PushMsg{
							Err:false,
							Code:200,
							Message:msg.Message,
						}
						send,_:=json.Marshal(newSend)
						client.send<-send
						is_read=1
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
			case "getHistory":
				c.getHistory(msg)
		}
	}
}
func checkMsgErr(err error) {
    if err != nil {
		fmt.Print(err)
		//TODO: 完善错误日志记录
    }
}
//TODO:待修改接口访问方式
func (c *Client) getHistory(msg PullMsg){
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/chat")
	if err != nil {
		fmt.Print(err)
		return 
	}
	defer func(){
		db.Close()
	}()
	rows,err:=db.Query("select uid,touid,msg from msg where is_read=? and ((uid=? and touid=?) or (uid=? and touid=?))",1,msg.Uuid,msg.ToUuid,msg.ToUuid,msg.Uuid)
	if err != nil {
		fmt.Print(err)
	}
	//TODO: 待修改历史数据格式以及传输方式
	for rows.Next(){
		sendMsg:=new(PushMsg)
		err = rows.Scan(&sendMsg.Uuid, &sendMsg.ToUuid,&sendMsg.Message)
		sendMsg.Err=false
		sendMsg.Code=200
		send,_:=json.Marshal(*sendMsg)
		c.send<-send
	}
}