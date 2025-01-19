// websocket/client.go
package websocket

import (
    "log"

    "github.com/gorilla/websocket"
)

// Client represents a WebSocket client
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

// NewClient creates a new WebSocket client instance
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
    return &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
    }
}

// ReadPump listens for messages from the WebSocket connection
func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    for {
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("[ERROR] Unexpected WebSocket close error: %v", err)
            } else {
                log.Printf("[INFO] WebSocket closed: %v", err)
            }
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
            err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
            if err != nil {
                log.Printf("[ERROR] WebSocket close message error: %v", err)
            }
            return
        }
        err := c.conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Printf("[ERROR] WebSocket write error: %v", err)
            return
        }
    }
}
