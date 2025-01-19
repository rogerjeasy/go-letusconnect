package handlers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/gofiber/websocket/v2"
// 	"github.com/rogerjeasy/go-letusconnect/models"
// 	"github.com/rogerjeasy/go-letusconnect/services"
// )

// type WebSocketHandler struct {
// 	wsService *services.WebSocketService
// }

// func NewWebSocketHandler(wsService *services.WebSocketService) *WebSocketHandler {
// 	return &WebSocketHandler{
// 		wsService: wsService,
// 	}
// }

// func (h *WebSocketHandler) HandleWebSocket(c *websocket.Conn) error {
// 	if c == nil {
// 		return fmt.Errorf("websocket connection cannot be nil")
// 	}

// 	// Get user ID from context (set by auth middleware)
// 	userID, ok := c.Locals("userID").(string)
// 	if !ok {
// 		return fmt.Errorf("user ID not found in context")
// 	}

// 	// Register connection
// 	h.wsService.RegisterConnection(userID, c)
// 	defer h.wsService.UnregisterConnection(userID)

// 	for {
// 		_, message, err := c.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("WebSocket read error: %v", err)
// 			}
// 			return fmt.Errorf("error reading message: %w", err)
// 		}

// 		var wsMessage models.WebSocketMessage
// 		if err := json.Unmarshal(message, &wsMessage); err != nil {
// 			log.Printf("Error unmarshaling message: %v", err)
// 			continue
// 		}

// 		wsMessage.From = userID
// 		wsMessage.Time = time.Now().Unix()

// 		// Handle different message types
// 		switch wsMessage.Type {
// 		case "direct_message":
// 			if err := h.wsService.HandleDirectMessage(wsMessage); err != nil {
// 				log.Printf("Error handling direct message: %v", err)
// 			}
// 		case "connection_request":
// 			if err := h.wsService.HandleConnectionRequest(wsMessage); err != nil {
// 				log.Printf("Error handling connection request: %v", err)
// 			}
// 		case "notification":
// 			if err := h.wsService.HandleNotification(wsMessage); err != nil {
// 				log.Printf("Error handling notification: %v", err)
// 			}
// 		default:
// 			log.Printf("Unknown message type: %s", wsMessage.Type)
// 		}
// 	}
// }
