package websocket

import (
	"time"
)

// Message represents a WebSocket message to be broadcast to clients.
type Message struct {
	Type      string    `json:"type"`      // e.g., "status_update"
	ID        string    `json:"id"`        // Identifier for the entity (e.g., pipeline run ID)
	Payload   string    `json:"payload"`   // JSON-encoded payload (e.g., status, message)
	Timestamp time.Time `json:"timestamp"` // Time of the message
}