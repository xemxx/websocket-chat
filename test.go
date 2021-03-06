package main

import(
	"time"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
    "net/http"
)


func main() {
    fmt.Println(time.Now().Unix())
    testz()
}

func for_switch(){
    for{
        fmt.Print("继续下一次")
        time.Sleep(time.Second*3)
        switch{
            case true:
                if(true){
                    continue
                }
        }
        fmt.Print("跳出switch")
        break
    }
    fmt.Print("跳出for")
}
type PushMsg struct{
	Err  bool		//是否错误
	Code int		//错误代码
	Uuid string 	//发送者
	ToUuid string	//接受者
	Message string  //具体数据
}
func inpro(){
    db,err:=sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/chat")
				if err != nil {
					fmt.Print(err)
                    db.Close()
                    return
				}
				rows,err:=db.Query("select uid,touid,msg from msg where is_read=?",0)
				if err != nil {
					fmt.Print(err)
                    db.Close()
                    return
				}
				for rows.Next(){
					sendMsg:=new(PushMsg)
					err = rows.Scan(&sendMsg.Uuid, &sendMsg.ToUuid,&sendMsg.Message)
					sendMsg.Err=false
                    sendMsg.Code=200
                    fmt.Println(sendMsg)
                    send,_:=json.Marshal(*sendMsg)
                    fmt.Println(string(send))
                }
				db.Close()
}

func httptest(){
    http.HandleFunc("/",handleHttp)
    err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func handleHttp(w http.ResponseWriter, r *http.Request){
    r.ParseForm()
    fmt.Println(r.PostForm)
    fmt.Println(r.PostFormValue("ad"))
}

type Msg struct {
	Uuid string 
	ToUuid string
	Msg string
	Send_time string
}

type PushMs struct{
	Err  bool		//是否错误
	Code int		//错误代码
	Uuid string 	//发送者
	ToUuid string	//接受者
	Data map[string]Msg  //具体数据
}
func testz(){
    push:=new(PushMs)
    push.Data=make(map[string]Msg)
	for i:=0;i<5;i++{
		msg:=new(Msg)
        msg.Msg="132"
        msg.Send_time="12345678900"
        msg.Uuid="1"
        msg.ToUuid="2"
        fmt.Println(msg)
        sty:=fmt.Sprintf("%d", i)
		push.Data[sty]=*msg
    }
    push.Data=nil
    send,_:=json.Marshal(*push)
    fmt.Println(string(send))

}