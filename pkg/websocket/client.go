package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // You might want to make this more restrictive in production
	},
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

// UpgradeProgressWS returns an Echo handler for progress notifications
func UpgradeProgressWS() echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		client := NewClient(ws)
		defer func() {
			GetHub().RemoveNotificationClient(client)
			ws.Close()
		}()

		// Register with hub
		GetHub().SetNotificationClient(client)

		// Set connection parameters
		ws.SetReadLimit(maxMessageSize)
		ws.SetReadDeadline(time.Now().Add(pongWait))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Keep connection alive
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}
		}

		return nil
	}
}

// UpgradeInteractiveWS returns an Echo handler for interactive communications
func UpgradeInteractiveWS() echo.HandlerFunc {
	return func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		client := NewClient(ws)
		defer func() {
			GetHub().RemoveInteractiveClient(client)
			ws.Close()
		}()

		// Register with hub
		GetHub().SetInteractiveClient(client)

		// Set connection parameters
		ws.SetReadLimit(maxMessageSize)
		ws.SetReadDeadline(time.Now().Add(pongWait))
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

		// Keep connection alive and handle messages
		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			// Only handle text messages
			if messageType == websocket.TextMessage {
				var response ActionResponse
				if err := json.Unmarshal(message, &response); err != nil {
					log.Printf("error unmarshaling message: %v", err)
					continue
				}

				// Send to hub for handling
				GetHub().HandleResponse(response)
			}
		}

		return nil
	}
}
