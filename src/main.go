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
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
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

