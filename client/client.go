package client

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	// "github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
	Db "websocket-chat/database"

)
var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
			return
		}
		//TODO: 解析msg的json finish
		go c.sovelMsg(msgJson)
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

func (client *Client)saveMsg(msg *PullMsg,is_read *int){
	mysql,err:= Db.NewMysql()
	if err !=nil {
		log.Fatal(err)
	}
	defer func() {
		mysql.Close()
	}()
	stmt,err:=mysql.Prepare("insert into msg(uid,touid,send_time,is_read,msg)values(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	res,err:=stmt.Exec(client.uuid,msg.ToUuid,time.Now().Unix(),is_read,msg.Message)
	if err != nil {
		log.Fatal(err)
	}
	msgId,_:=res.LastInsertId()
	a, _ := strconv.ParseInt(msg.Uuid, 10, 64)
	b, _ :=strconv.ParseInt(msg.ToUuid, 10, 64)
	c:=msg.Uuid
	d:=msg.ToUuid
	if a<b {
		e:=c
		c=d
		d=e
	}
	sql:="insert into msglist (uid,touid,msg_id,num)values(?,?,?,1) ON DUPLICATE KEY UPDATE msg_id=?,num=num+1"
	if *is_read == 1{
		sql="insert into msglist (uid,touid,msg_id,num)values(?,?,?,0) ON DUPLICATE KEY UPDATE msg_id=?,num=0"
	}
	rows,err:=mysql.Query(sql,c,d,msgId,msgId)
	if err != nil {
		fmt.Print(err)
	}
	defer rows.Close()
}

func (c *Client)sovelMsg(msgJson []byte){
	mysql,err:= Db.NewMysql()
	if err !=nil {
		log.Fatal(err)
	}
	defer func() {
		mysql.Close()
	}()
	msg:= &PullMsg{}
	if json.Unmarshal(msgJson,msg) != nil{
		//TODO:返回错误信息  finish 
		c.send<-sendMsg(&PushMsg{"",true,401,"json error","",""})
		return
	}
	if (msg.Uuid!=c.uuid) && c.isBind(){
		c.send<-sendMsg(&PushMsg{"",true,405,"uid error","",""})
		return
	}
	if !c.isBind() && msg.Type != "bind" {
		c.send<-sendMsg(&PushMsg{msg.Type,true,402,"请先绑定后请求","",""})
		return
	}
	//TODO: 通过账户寻找发送账户是否在线并推送  finish
	switch msg.Type{
		case "bind":
			if c.isBind(){
				c.send<-sendMsg(&PushMsg{msg.Type,true,403,"请勿重复绑定","",""})
				return
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
			is_read:=0
			for client:=range manager.clients{
				//判断是否在线
				if msg.ToUuid == client.uuid{
					if client.join == msg.Uuid{
						client.send<-sendMsg(&PushMsg{msg.Type,false,200,msg.Message,msg.Uuid,""})
						is_read=1
					}
					break
				}
			}
			go c.saveMsg(msg,&is_read)

		case "join":
			c.join=msg.ToUuid
			c.send<-sendMsg(&PushMsg{msg.Type,false,202,"join success","",""})

			a, _ := strconv.ParseInt(msg.Uuid, 10, 64)
			b, _ :=strconv.ParseInt(msg.ToUuid, 10, 64)
			c:=msg.Uuid
			d:=msg.ToUuid
			if a<b {
				e:=c
				c=d
				d=e
			}
			rows,err:=mysql.Query("insert ignore into msglist (uid,touid)values(?,?)",c,d)
			if err != nil {
				fmt.Print(err)
				return
			}
			rows.Close()

		case "exit":
			c.join=""
			c.send<-sendMsg(&PushMsg{msg.Type,false,200,"exit success","",""})
	}
}
