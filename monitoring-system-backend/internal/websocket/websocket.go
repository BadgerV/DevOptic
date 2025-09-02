package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Hub manages WebSocket connections and broadcasts messages.
type Hub struct {
	logger     *zap.Logger
	clients    map[string]map[*websocket.Conn]bool // Map of entity ID to clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mu         sync.RWMutex
	upgrader   websocket.Upgrader
}

// Client represents a single WebSocket connection for an entity ID.
type Client struct {
	conn     *websocket.Conn
	entityID string
}

// NewHub creates a new WebSocket Hub.
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		logger:     logger,
		clients:    make(map[string]map[*websocket.Conn]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true }, // Adjust for production
		},
	}
}

// Run starts the WebSocket hub, handling client registration, unregistration, and broadcasts.
func (h *Hub) Run(ctx context.Context) {
	fmt.Println("Getting here")
	for {
		select {
		case <-ctx.Done():
			h.logger.Info("Shutting down WebSocket hub")
			fmt.Println("Shutting down WebSocket hub")

			return
		case client := <-h.register:
			h.mu.Lock()
			if _, exists := h.clients[client.entityID]; !exists {
				h.clients[client.entityID] = make(map[*websocket.Conn]bool)
			}
			h.clients[client.entityID][client.conn] = true
			clientCount := len(h.clients[client.entityID])
			h.mu.Unlock()
			h.logger.Info("Client registered", zap.String("entity_id", client.entityID))
			fmt.Printf("Client registered for entity_id: %s (total clients for this entity: %d)\n", client.entityID, clientCount)
		case client := <-h.unregister:
			h.mu.Lock()
			if clients, exists := h.clients[client.entityID]; exists {
				delete(clients, client.conn)
				if len(clients) == 0 {
					delete(h.clients, client.entityID)
				}
			}
			h.mu.Unlock()
			client.conn.Close()
			h.logger.Info("Client unregistered", zap.String("entity_id", client.entityID))
			fmt.Printf("Client unregistered for entity_id: %s\n", client.entityID)

		case message := <-h.broadcast:
			fmt.Printf("[BROADCAST] Attempting to broadcast message to entity_id: %s, type: %s\n", message.ID, message.Type)
			
			h.mu.RLock()
			clients, exists := h.clients[message.ID]
			if !exists {
				h.mu.RUnlock()
				fmt.Printf("[BROADCAST] No clients found for entity_id: %s - broadcast skipped\n", message.ID)
				continue
			}
			
			clientCount := len(clients)
			fmt.Printf("[BROADCAST] Found %d client(s) for entity_id: %s\n", clientCount, message.ID)
			
			data, err := json.Marshal(message)
			if err != nil {
				h.logger.Error("Failed to marshal broadcast message", zap.Error(err))
				fmt.Printf("[BROADCAST] Failed to marshal broadcast message for entity_id: %s - %v\n", message.ID, err)
				h.mu.RUnlock()
				continue
			}
			
			successCount := 0
			failCount := 0
			
			for conn := range clients {
				err := conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					h.logger.Error("Failed to send WebSocket message", zap.String("entity_id", message.ID), zap.Error(err))
					fmt.Printf("[BROADCAST] Failed to send message to client for entity_id: %s - %v\n", message.ID, err)
					failCount++
					h.mu.RUnlock()
					h.unregister <- &Client{conn: conn, entityID: message.ID}
					continue
				}
				successCount++
			}
			h.mu.RUnlock()
			
			fmt.Printf("[BROADCAST] Broadcast completed for entity_id: %s, type: %s - Success: %d, Failed: %d\n", 
				message.ID, message.Type, successCount, failCount)
			
			h.logger.Info("Broadcast message sent", 
				zap.String("entity_id", message.ID), 
				zap.String("type", message.Type),
				zap.Int("success_count", successCount),
				zap.Int("fail_count", failCount))
		}
	}
}

// HandleWebSocket upgrades an HTTP connection to a WebSocket and registers the client.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, entityID string) {
	fmt.Printf("HandleWebSocket called for entity_id: %s\n", entityID)
	fmt.Printf("Request headers: %+v\n", r.Header)
	
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed for entity_id: %s - Error: %v\n", entityID, err)
		h.logger.Error("Failed to upgrade WebSocket connection", zap.String("entity_id", entityID), zap.Error(err))
		http.Error(w, "Failed to upgrade WebSocket connection", http.StatusInternalServerError)
		return
	}
	
	fmt.Printf("WebSocket upgrade successful for entity_id: %s\n", entityID)

	client := &Client{
		conn:     conn,
		entityID: entityID,
	}
	
	fmt.Printf("Client object created, sending to register channel for entity_id: %s\n", entityID)
	h.register <- client
	fmt.Printf("Client sent to register channel for entity_id: %s\n", entityID)

	// Handle incoming messages (optional, for client-initiated requests)
	go func() {
		defer func() {
			h.unregister <- client
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				h.logger.Error("WebSocket read error", zap.String("entity_id", entityID), zap.Error(err))
				return
			}
			// Handle client messages if needed (e.g., ping/pong)
		}
	}()
}

// Broadcast sends a message to all clients subscribed to an entity ID.
func (h *Hub) Broadcast(message Message) {
	fmt.Printf("[BROADCAST] Queuing message for broadcast - entity_id: %s, type: %s\n", message.ID, message.Type)
	h.broadcast <- message
}