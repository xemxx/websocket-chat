package main

import (
	"net/http"
	"log"
	"flag"
	"websocket-chat/src/client"
)

var addr = flag.String("addr", ":8080", "http service address")

func main(){
	flag.Parse()
	go client.Run()
	http.HandleFunc("/",handleHttp)
	http.HandleFunc("/ws", client.HandleWs)
	http.HandleFunc("/login", client.Login)
	http.HandleFunc("/logout", client.Logout)
	http.HandleFunc("/user/gethistory",handleGetHistory)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


