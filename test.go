package main

import(
	"time"
	"fmt"
)


func main() {
    fmt.Print(time.Now().Unix())
    for_switch()
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