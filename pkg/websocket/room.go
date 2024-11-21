package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	name       string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan *Message
	mu         sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *Message),
	}
}

func (room *Room) Start() {
	for {
		select {
		case client := <-room.Register:
			room.mu.Lock()
			room.Clients[client] = true
			room.mu.Unlock()

			fmt.Println("Size of Connection Pool: ", len(room.Clients))

			for client := range room.Clients {
				fmt.Println(client)
				notify := &Message{ClientName: client.ID, Text: "New User Joined..."}
				err := client.Conn.WriteMessage(websocket.TextMessage, getTemplate("templates/notify.html", notify))
				if err != nil {
					fmt.Println("Error writing notify", err)
					continue
				}
			}

		case client := <-room.Unregister:
			room.mu.Lock()
			delete(room.Clients, client)
			room.mu.Unlock()
			client.Conn.Close()

			fmt.Println("Size of Connection Pool: ", len(room.Clients))

			for client := range room.Clients {
				notify := &Message{ClientName: client.ID, Text: "User Disconnected..."}
				err := client.Conn.WriteMessage(websocket.TextMessage, getTemplate("templates/notify.html", notify))
				if err != nil {
					fmt.Println("Error writing notify", err)
					continue
				}
			}

		case message := <-room.Broadcast:
			room.mu.Lock()
			for client := range room.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, getTemplate("templates/message.html", message))
				if err != nil {
					fmt.Println("Error writing message", err)
					continue
				}
			}
			room.mu.Unlock()
		}
	}
}
