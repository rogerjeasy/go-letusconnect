// websocket/client.go
package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/rogerjeasy/go-letusconnect/models"
)

func NewClient(conn *websocket.Conn, manager *models.Manager, userID string) *models.Client {
	return &models.Client{
		Conn:     conn,
		Manager:  manager,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		IsActive: true,
	}
}

func ReadPump(c *models.Client) {
	ctx := context.Background()

	defer func() {
		c.Manager.Unregister <- c
		log.Printf("Client disconnected: %s", c.UserID)
		c.Conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		messageType, message, err := c.Conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Printf("Client closed connection normally: %s", c.UserID)
				return
			}
			log.Printf("Read error (client %s): %v", c.UserID, err)
			return
		}

		if messageType != websocket.MessageText {
			continue
		}

		var msg models.WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message (client %s): %v", c.UserID, err)
			continue
		}

		msg.From = c.UserID
		messageBytes, _ := json.Marshal(msg)
		c.Manager.Broadcast <- messageBytes
	}
}

func WritePump(c *models.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ensure the connection is closed properly
	defer func() {
		log.Printf("Closing connection for client: %s", c.UserID)
		c.Conn.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	// Heartbeat ticker for keeping the connection alive
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			// If the channel is closed, terminate the connection
			if !ok {
				log.Printf("Channel closed for client: %s", c.UserID)
				c.Conn.Close(websocket.StatusNormalClosure, "channel closed")
				return
			}

			// Write the message to the WebSocket connection
			if err := c.Conn.Write(ctx, websocket.MessageText, message); err != nil {
				log.Printf("Write error for client %s: %v", c.UserID, err)
				return
			}

		case <-ticker.C:
			// Send a ping message to keep the connection alive
			if err := c.Conn.Ping(ctx); err != nil {
				log.Printf("Ping failed for client %s: %v", c.UserID, err)
				return
			}
		}
	}
}
