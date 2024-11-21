package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan *Message
	mu         sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.mu.Lock()
			pool.Clients[client] = true
			pool.mu.Unlock()

			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			for client := range pool.Clients {
				fmt.Println(client)
				notify := &Message{ClientName: client.ID, Text: "New User Joined..."}
				err := client.Conn.WriteMessage(websocket.TextMessage, getTemplate("templates/notify.html", notify))
				if err != nil {
					fmt.Println("Error writing notify", err)
					continue
				}
			}

		case client := <-pool.Unregister:
			pool.mu.Lock()
			delete(pool.Clients, client)
			pool.mu.Unlock()
			client.Conn.Close()

			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			for client := range pool.Clients {
				notify := &Message{ClientName: client.ID, Text: "User Disconnected..."}
				err := client.Conn.WriteMessage(websocket.TextMessage, getTemplate("templates/notify.html", notify))
				if err != nil {
					fmt.Println("Error writing notify", err)
					continue
				}
			}

		case message := <-pool.Broadcast:
			pool.mu.Lock()
			for client := range pool.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, getTemplate("templates/message.html", message))
				if err != nil {
					fmt.Println("Error writing message", err)
					continue
				}
			}
			pool.mu.Unlock()
		}
	}
}
