package handlers

import (
	"context"
	"sort"
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

	// Create a sorted list of sender and receiver IDs
	ids := []string{message.SenderID, message.ReceiverID}
	sort.Strings(ids)
	channelName := "private-messages-" + strings.Join(ids, "-")

	// Trigger Pusher event with the consistent channel name
	err = services.PusherClient.Trigger(
		channelName,
		"new-message",
		mappers.MapMessageGoToFrontend(message),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger event",
		})
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

// SendTyping handles the typing event and triggers a Pusher event
func SendTyping(c *fiber.Ctx) error {
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

	// Parse request payload
	var payload struct {
		ReceiverID string `json:"receiverId"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Ensure receiver ID is provided
	if payload.ReceiverID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "receiverId is required",
		})
	}

	// Create a sorted list of sender and receiver IDs
	ids := []string{uid, payload.ReceiverID}
	sort.Strings(ids)
	channelName := "private-messages-" + strings.Join(ids, "-")

	// Trigger Pusher event with the consistent channel name
	err = services.PusherClient.Trigger(channelName, "user-typing", map[string]string{
		"senderId":   uid,
		"receiverId": payload.ReceiverID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send typing notification",
		})
	}

	return c.JSON(fiber.Map{
		"success": "Typing notification sent",
	})
}

// SendDirectMessage handles sending a direct message and triggering a Pusher event
func SendDirectMessage(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload."})
	}

	message := mappers.MapDirectMessageFrontendToGo(payload)
	message.ID = uuid.New().String()

	if message.SenderID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only send messages as yourself."})
	}

	ids := []string{message.SenderID, message.ReceiverID}
	sort.Strings(ids)
	channelID := strings.Join(ids, "-")

	docRef := services.FirestoreClient.Collection("messages").Doc(channelID)

	docSnapshot, err := docRef.Get(context.Background())
	if err != nil && !docSnapshot.Exists() {
		messages := models.Messages{
			ChannelID:      channelID,
			DirectMessages: []models.DirectMessage{message},
		}

		_, err = docRef.Set(context.Background(), mappers.MapMessagesGoToFirestore(messages))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message."})
		}
	} else {
		var existingMessages models.Messages
		err := docSnapshot.DataTo(&existingMessages)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read existing messages."})
		}

		existingMessages.DirectMessages = append(existingMessages.DirectMessages, message)

		_, err = docRef.Set(context.Background(), mappers.MapMessagesGoToFirestore(existingMessages))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update messages."})
		}
	}

	// Trigger Pusher event with the consistent channel name
	channelName := "private-messages-" + channelID
	err = services.PusherClient.Trigger(
		channelName,
		"new-direct-message",
		mappers.MapDirectMessageGoToFrontend(message),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger event.",
		})
	}

	return c.JSON(fiber.Map{"success": "Direct message sent successfully.", "message": message})
}

// GetDirectMessages fetches direct messages based on the channel ID
func GetDirectMessages(c *fiber.Ctx) error {

	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	// Validate token and get UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}
	// Get sender and receiver IDs from query parameters
	senderID := c.Query("senderId")
	receiverID := c.Query("receiverId")

	if senderID == "" || receiverID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both senderId and receiverId are required.",
		})
	}

	// Create a sorted list of sender and receiver IDs for consistent channel ID
	ids := []string{senderID, receiverID}
	sort.Strings(ids)
	channelID := strings.Join(ids, "-")

	// Reference to the Firestore document
	docRef := services.FirestoreClient.Collection("messages").Doc(channelID)

	// Fetch the document from Firestore
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil || !docSnapshot.Exists() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success":  "No messages found, starting a new conversation.",
			"messages": []models.DirectMessage{},
		})
	}

	// Convert Firestore data to the Messages struct
	var messages models.Messages
	err = docSnapshot.DataTo(&messages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse messages.",
		})
	}

	// Convert messages to frontend format
	frontendMessages := mappers.MapMessagesGoToFrontend(messages)

	return c.JSON(fiber.Map{
		"success":  "Messages fetched successfully.",
		"messages": frontendMessages,
	})
}

// ========================= Send Group Message =========================

// SendGroupMessage handles sending a group message and triggering a Pusher event
func SendGroupMessage(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload."})
	}

	// Convert payload to GroupMessage struct
	message := mappers.MapGroupMessageFrontendToGo(payload)
	message.ID = uuid.New().String()

	// Ensure the sender ID matches the authenticated user
	if message.SenderID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You can only send messages as yourself."})
	}

	// Add message to Firestore
	_, _, err = services.FirestoreClient.Collection("group_messages").Add(context.Background(), mappers.MapGroupMessageGoToFirestore(message))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send group message."})
	}

	// Create a channel name based on the project ID or group ID
	channelName := "group-messages-" + message.ProjectID
	if message.GroupID != nil {
		channelName += "-" + *message.GroupID
	}

	// Trigger Pusher event with the group channel name
	err = services.PusherClient.Trigger(
		channelName,
		"new-group-message",
		mappers.MapGroupMessageGoToFrontend(message),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger event.",
		})
	}

	return c.JSON(fiber.Map{"success": "Group message sent successfully.", "message": message})
}
