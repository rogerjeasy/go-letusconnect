package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// SendMessage handles sending a message and triggering a Pusher event
func SendMessage(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again",
		})
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Convert payload to Message struct
	message := mappers.MapMessageFrontendToGo(payload)
	message.ID = uuid.New().String()
	message.CreatedAt = time.Now()

	// Ensure the sender ID matches the authenticated user
	if message.SenderID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only send messages as yourself"})
	}

	// Add message to Firestore
	_, _, err = services.FirestoreClient.Collection("messages").Add(context.Background(), mappers.MapMessageGoToFirestore(message))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message"})
	}

	// Trigger Pusher event for real-time notification
	err = services.PusherClient.Trigger("messages-channel", "new-message", mappers.MapMessageGoToFrontend(message))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to trigger event"})
	}

	return c.JSON(fiber.Map{"success": "Message sent successfully", "message": message})
}

// GetMessages retrieves messages between two users
func GetMessages(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	senderID := c.Query("senderId")
	receiverID := c.Query("receiverId")

	if senderID == "" || receiverID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "senderId and receiverId are required"})
	}

	// Ensure the authenticated user is either the sender or the receiver
	if uid != senderID && uid != receiverID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to view these messages"})
	}

	// Query Firestore for messages
	iter := services.FirestoreClient.Collection("messages").
		Where("sender_id", "in", []string{senderID, receiverID}).
		Where("receiver_id", "in", []string{senderID, receiverID}).
		Documents(context.Background())

	var messages []models.Message
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		msg := mappers.MapMessageFirestoreToGo(doc.Data())
		messages = append(messages, msg)
	}

	// Convert messages to frontend format
	frontendMessages := mappers.MapMessagesArrayToFrontend(messages)

	return c.JSON(frontendMessages)
}
