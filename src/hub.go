package main

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

func (h *ClientManager) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// case message := <-h.broadcast:
		// 	for client := range h.clients {
		// 		select {
		// 		case client.send <- message:
		// 		default:
		// 			close(client.send)
		// 			delete(h.clients, client)
		// 		}
		// 	}
		}
	}
}