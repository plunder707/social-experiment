// websocket/client.go
package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// ReadPump listens for messages from the WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[INFO] WebSocket closed: %v", err)
			break
		}
		// Handle incoming messages if necessary
	}
}

// WritePump sends messages to the WebSocket connection
func (c *Client) WritePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			// Hub closed the channel
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("[ERROR] WebSocket write error: %v", err)
			return
		}
	}
}
