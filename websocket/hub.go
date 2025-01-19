// websocket/hub.go
package websocket

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"social-experiment/models"

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
			log.Printf("[INFO] Client registered: %v", client.conn.RemoteAddr())
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("[INFO] Client unregistered: %v", client.conn.RemoteAddr())
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
					log.Printf("[WARNING] Client send channel full, removed client: %v", client.conn.RemoteAddr())
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
	// Extract and validate JWT token from Authorization header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Assuming the token is prefixed with "Bearer "
	parts := strings.SplitN(tokenString, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	tokenString = parts[1]

	// Parse and verify JWT token
	claims, err := ParseJWT(tokenString, "your_jwt_secret") // Replace with actual secret
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Upgrade to WebSocket
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

	client := NewClient(h, conn)
	h.register <- client

	// Start read and write pumps
	go client.ReadPump()
	go client.WritePump()

	log.Printf("[INFO] New WebSocket connection established: %v", conn.RemoteAddr())
}
