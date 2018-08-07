package main

import (
	"net/http"
	"log"
)

func main(){
	http.HandleFunc("/",handerHttp)
	http.HandleFunc("/ws", handerWs)
	
}

func handerHttp(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func handerWs(w http.ResponseWriter, r *http.Request){

}