package websocket

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (c *Client) SendProgress(progress float64, message string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	msg := ProgressMessage{
		Type:     "progress",
		Progress: progress,
		Message:  message,
	}

	return c.Conn.WriteJSON(msg)
}

func (c *Client) SendJSON(v interface{}) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return c.Conn.WriteJSON(v)
}

func (c *Client) ReadJSON(v interface{}) error {
	return c.Conn.ReadJSON(v)
}

// UpgradeProgressWS returns a Fiber middleware for progress notifications
func UpgradeProgressWS() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		client := NewClient(c)
		defer func() {
			GetHub().RemoveNotificationClient(client)
			c.Close()
		}()

		// Register with hub
		GetHub().SetNotificationClient(client)

		// Set connection parameters
		c.SetReadLimit(maxMessageSize)
		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Keep connection alive
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}
		}
	}, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
}

// UpgradeInteractiveWS returns a Fiber middleware for interactive communications
func UpgradeInteractiveWS() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		client := NewClient(c)
		defer func() {
			GetHub().RemoveInteractiveClient(client)
			c.Close()
		}()

		// Register with hub
		GetHub().SetInteractiveClient(client)

		// Set connection parameters
		c.SetReadLimit(maxMessageSize)
		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Keep connection alive and handle messages
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}
		}
	}, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
}
