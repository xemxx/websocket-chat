package main

import (
	"flag"
	"fmt"
	"net/http"
	// "websocket-chat/newhttp"
	"websocket-chat/client"
)

var addr = flag.String("addr", ":8080", "http service address")

func main(){
	flag.Parse()
	go client.Run()
	//http.HandleFunc("/",newhttp.HandleHttp)
	http.HandleFunc("/ws", client.HandleWs)
	// http.HandleFunc("/user/gethistory",newhttp.HandleGetHistory)
	fmt.Println("Listen ",*addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
	 	fmt.Println("ListenAndServe: ", err)
	}
}


