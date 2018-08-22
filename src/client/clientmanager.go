package client

// 聊天室类型
type ClientManager struct {
	clients map[*Client]bool  // 已注册的连接用户
	// broadcast chan []byte // 推送消息通道
	register chan *Client // 注册通道
	unregister chan *Client // 注销通道
}

func newClientManager() *ClientManager {
	return &ClientManager{
		//broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

var manager=newClientManager()


func Run() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client] = true
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
			}
		}
	}
}