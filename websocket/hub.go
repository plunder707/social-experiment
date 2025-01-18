// websocket/hub.go
package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"maliaki-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewHub initializes a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastPost sends a new post to all connected clients
func (h *Hub) BroadcastPost(post models.Post) {
	postJSON, err := json.Marshal(post)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal post: %v", err)
		return
	}
	h.broadcast <- postJSON
}

// HandleWebSocket handles incoming WebSocket connections
func (h *Hub) HandleWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Implement origin checking based on CORS if necessary
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ERROR] WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	h.register <- client

	// Start read and write pumps
	go client.ReadPump()
	go client.WritePump()
}
