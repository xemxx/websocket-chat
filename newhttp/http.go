package newhttp

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"websocket-chat/mysql"
)

type Msg struct {
	Uuid string 
	ToUuid string
	Msg string
	Send_time string
}

type PushMsg struct{
	Err  bool		//是否错误
	Code int		//错误代码
	Uuid string 	//发送者
	ToUuid string	//接受者
	Data map[string]Msg  //具体数据
	ErrMsg string	//具体错误消息
}

func HandleHttp(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../web/home.html")
}
//TODO:优化接口，准确推送阅读消息
func HandleGetHistory(w http.ResponseWriter, r *http.Request){
	log.Println(r.URL)
	db := mysql.NewMysql()
	defer db.Close()
	posts:=r.PostForm
	rows,err:=db.Query("select uid,touid,msg,send_time from msg where is_read=? and  ((uid=? and touid=?) or (uid=? and touid=?)) order by send_time desc limit ?",1,posts["uid"],posts["touid"],posts["touid"],posts["uid"],posts["num"])
	if err != nil {
		fmt.Println(err)
	}
	//TODO: 待修改历史数据格式以及传输方式
	defer rows.Close()
	push:=new(PushMsg)
	i:=0
	push.Data=make(map[string]Msg)
	for rows.Next(){
		msg:=new(Msg)
		err = rows.Scan(&msg.Uuid, &msg.ToUuid,&msg.Msg,&msg.Send_time)
		str:=fmt.Sprintf("%d", i)
		push.Data[str]=*msg
		i++
	}
	push.Err=false
	push.Code=200
	//由于http包的参数解析将同名参数解析为map，因此需要默认第一个
	push.Uuid=posts["uid"][0]
	push.ToUuid=posts["touid"][0]
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		push.Err=true
		push.Code=500
		push.ErrMsg="系统读写错误"
		push.Data=nil
	}
	send,_:=json.Marshal(*push)
	w.Header().Set("Content-Type", "application/json")
    w.Write(send)
}