package websocket

import (
	"sync"

	"neploy.dev/pkg/logger"
)

// Global hub instance
var globalHub = NewHub()

// Hub manages WebSocket clients
type Hub struct {
	notification *Client
	interactive  *Client
	mu           sync.Mutex
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{}
}

// SetNotificationClient sets the notification client
func (h *Hub) SetNotificationClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notification = client
}

// SetInteractiveClient sets the interactive client
func (h *Hub) SetInteractiveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.interactive = client
}

// RemoveNotificationClient removes the notification client if it matches
func (h *Hub) RemoveNotificationClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.notification == client {
		h.notification = nil
	}
}

// RemoveInteractiveClient removes the interactive client if it matches
func (h *Hub) RemoveInteractiveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.interactive == client {
		h.interactive = nil
	}
}

// GetNotificationClient gets the notification client
func (h *Hub) GetNotificationClient() *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.notification // can be nil
}

// BroadcastProgress sends a progress message to the notification client
func (h *Hub) BroadcastProgress(progress float64, message string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.notification == nil {
		return
	}

	err := h.notification.SendProgress(progress, message)
	if err != nil {
		logger.Error("error sending progress: %v", err)
	}
}

// BroadcastInteractive sends an action message to the interactive client
func (h *Hub) BroadcastInteractive(msg ActionMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.interactive == nil {
		return
	}

	err := h.interactive.SendJSON(msg)
	if err != nil {
		logger.Error("error sending interactive message: %v", err)
	}
}

// GetHub returns the global hub instance
func GetHub() *Hub {
	return globalHub
}
