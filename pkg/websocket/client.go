package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

type ProgressMessage struct {
	Type     string  `json:"type"`
	Progress float64 `json:"progress"`
	Message  string  `json:"message"`
}

func NewClient(c *websocket.Conn) *Client {
	return &Client{
		Conn: c,
	}
}

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
		defer c.Close()

		// Register with hub
		GetHub().SetNotificationsClient(client)

		// Set connection parameters
		c.SetReadLimit(maxMessageSize)
		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Start ping ticker
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		// Message handling loop - for progress, we only send updates
		for {
			select {
			case <-ticker.C:
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("error sending ping: %v", err)
					return
				}
			default:
				// For progress WS, we only handle connection maintenance
				_, _, err := client.Conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					return
				}
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
		defer c.Close()

		// Register with hub
		GetHub().SetInteractiveClient(client)

		// Set connection parameters
		c.SetReadLimit(maxMessageSize)
		c.SetReadDeadline(time.Now().Add(pongWait))
		c.SetPongHandler(func(string) error {
			c.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Start ping ticker
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		// Message handling loop - for interactive, we handle both read and write
		for {
			select {
			case <-ticker.C:
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("error sending ping: %v", err)
					return
				}
			default:
				messageType, message, err := client.Conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					return
				}

				switch messageType {
				case websocket.TextMessage:
					// Handle interactive messages
					if err := client.Conn.WriteMessage(messageType, message); err != nil {
						log.Printf("error echo message: %v", err)
						return
					}
				case websocket.PingMessage:
					if err := client.Conn.WriteMessage(websocket.PongMessage, nil); err != nil {
						log.Printf("error sending pong: %v", err)
						return
					}
				}
			}
		}
	}, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
}
