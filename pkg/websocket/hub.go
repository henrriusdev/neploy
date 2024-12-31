package websocket

import (
	"sync"
	"time"

	"neploy.dev/pkg/logger"
)

// Global hub instance
var globalHub = NewHub()

// Hub manages WebSocket clients
type Hub struct {
	notification *Client
	interactive  *Client
	mu           sync.Mutex
	responseCh   chan ActionResponse
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		responseCh: make(chan ActionResponse, 1),
	}
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

// GetInteractiveClient gets the interactive client
func (h *Hub) GetInteractiveClient() *Client {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.interactive // can be nil
}

// HandleResponse handles a response from a client
func (h *Hub) HandleResponse(response ActionResponse) {
	select {
	case h.responseCh <- response:
		logger.Info("response queued: %+v", response)
	default:
		logger.Error("response channel full, dropping: %+v", response)
	}
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

// BroadcastInteractive sends an action message to the interactive client and waits for response
func (h *Hub) BroadcastInteractive(msg ActionMessage) *ActionResponse {
	h.mu.Lock()
	if h.interactive == nil {
		h.mu.Unlock()
		logger.Info("no interactive client connected")
		return nil
	}

	err := h.interactive.SendJSON(msg)
	h.mu.Unlock()

	if err != nil {
		logger.Error("error sending interactive message: %v", err)
		return nil
	}

	// Wait for response with timeout
	select {
	case response := <-h.responseCh:
		logger.Info("received response from interactive client: %+v", response)
		return &response
	case <-time.After(30 * time.Second):
		logger.Error("timeout waiting for response")
		return nil
	}
}

// GetHub returns the global hub instance
func GetHub() *Hub {
	return globalHub
}
