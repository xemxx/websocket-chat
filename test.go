package main

import(
	"time"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "encoding/json"
)


func main() {
    fmt.Print(time.Now().Unix())
    inpro()
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