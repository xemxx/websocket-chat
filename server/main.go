package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xemxx/websocket-chat/server/client"
	"go.uber.org/zap"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	go client.Run()
	l, _ := zap.NewProduction()
	zap.ReplaceGlobals(l)

	router := gin.Default()
	router.GET("/ws", client.HandleWs)
	router.POST("/login", client.Login)
	router.GET("/logout", client.Logout)
	fmt.Println("Listen ", *addr)
	err := http.ListenAndServe(*addr, router)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
