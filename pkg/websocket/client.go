package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Fatal("Error Unmarshalling")
		}

		c.Pool.Broadcast <- &message
	}
}
