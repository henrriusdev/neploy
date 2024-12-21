package websocket

import (
	"sync"
)

var (
	// Global hub instance
	globalHub = NewHub()
)

// Hub manages WebSocket clients
type Hub struct {
	notificationsClient *Client
	interactiveClient  *Client
	mu                 sync.Mutex
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{}
}

// SetNotificationsClient sets the notifications client
func (h *Hub) SetNotificationsClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notificationsClient = client
}

// SetInteractiveClient sets the interactive client
func (h *Hub) SetInteractiveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.interactiveClient = client
}

// GetNotificationsClient gets the notifications client
func (h *Hub) GetNotificationsClient() *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.notificationsClient
}

// GetInteractiveClient gets the interactive client
func (h *Hub) GetInteractiveClient() *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.interactiveClient
}

// GetHub returns the global hub instance
func GetHub() *Hub {
	return globalHub
}
