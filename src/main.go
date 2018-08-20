package main

import (
	"net/http"
	"log"
	"flag"
)

var addr = flag.String("addr", ":8080", "http service address")

func main(){
	flag.Parse()
	manager := newClientManager()
	go manager.run()
	http.HandleFunc("/",handleHttp)
	http.HandleFunc("/ws", func(w http.ResponseWriter,r *http.Request){
		HandleWs(manager,w,r)
	})
	//http.HandleFunc("/ws", HandleWs)
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
