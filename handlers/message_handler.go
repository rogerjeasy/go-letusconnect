package handlers

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
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

// GetMessages retrieves all messages for the authenticated user
func GetMessages(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	iter := services.FirestoreClient.Collection("messages").
		Documents(context.Background())

	var directMessages []models.DirectMessage
	var channelID string

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to fetch messages: %v", err),
			})
		}

		data := doc.Data()

		channelID, ok := data["channel_id"].(string)
		if ok && strings.Contains(channelID, uid) {
			directMessage := mappers.MapDirectMessageFirestoreToGo(data)
			fmt.Println("directMessages: ", directMessage)
			directMessages = append(directMessages, directMessage)
		}
	}

	if len(directMessages) == 0 {
		directMessages = []models.DirectMessage{}
		channelID = ""
	}

	messages := models.Messages{
		ChannelID:      channelID,
		DirectMessages: directMessages,
	}

	frontendMessages := mappers.MapMessagesGoToFrontend(messages)

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

	receiverUserUid, ok := payload["receiverId"].(string)
	if !ok || receiverUserUid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "receiverId is required."})
	}

	receiverUser, err := services.GetUserByUID(receiverUserUid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch receiver details",
		})
	}
	receiverName := receiverUser["username"].(string)
	createdAt := time.Now().Format(time.RFC3339)

	fmt.Println("receiverName: ", receiverName)

	// now add receiverName to the payload
	payload["receiverName"] = receiverName
	payload["createdAt"] = createdAt

	message := mappers.MapDirectMessageFrontendToGo(payload)
	message.ID = uuid.New().String()
	message.ReadStatus = map[string]bool{message.SenderID: true, message.ReceiverID: false}

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

	// Trigger Pusher event for notification
	notificationChannel := "user-notifications-direct-msg-" + message.ReceiverID
	err = services.PusherClient.Trigger(
		notificationChannel,
		"update-unread-count",
		map[string]string{
			"senderName": message.SenderName,
			"content":    message.Content,
			"senderID":   message.SenderID,
			"receiverId": message.ReceiverID,
			"messageId":  message.ID,
		},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger notification event.",
		})
	}

	return c.JSON(fiber.Map{"success": "Direct message sent successfully.", "message": message})
}

// GetDirectMessages fetches direct messages for the authenticated user based on their UID
func GetDirectMessages(c *fiber.Ctx) error {

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

	// Fetch all documents from the Firestore messages collection
	iter := services.FirestoreClient.Collection("messages").Documents(context.Background())

	var messagesList []models.Messages

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch messages.",
			})
		}

		data := doc.Data()
		channelID, ok := data["channel_id"].(string)
		if ok && strings.Contains(channelID, uid) {
			var messages models.Messages
			err = doc.DataTo(&messages)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to parse messages.",
				})
			}
			messagesList = append(messagesList, messages)
		}
	}

	// Check if no messages were found
	if len(messagesList) == 0 {
		return c.JSON(fiber.Map{
			"success":  "No messages found.",
			"messages": []map[string]interface{}{},
		})
	}

	// Convert all messages to frontend format
	var frontendMessages []map[string]interface{}
	for _, msg := range messagesList {
		frontendMessages = append(frontendMessages, mappers.MapMessagesGoToFrontend(msg))
	}

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

	// Fetch the list of group members (assuming there's a function to get group members)
	groupMembers, err := services.GetGroupMembers(message.ProjectID, message.GroupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch group members.",
		})
	}

	// Initialize ReadStatus for all group members
	message.ReadStatus = make(map[string]bool)
	for _, memberID := range groupMembers {
		// Set the sender's ReadStatus to true and others to false
		if memberID == message.SenderID {
			message.ReadStatus[memberID] = true
		} else {
			message.ReadStatus[memberID] = false
		}
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

// GetUnreadMessagesCount fetches the count of unread direct messages for the logged-in user.
// If a senderId query parameter is provided, it counts unread messages only from that sender.
func GetUnreadMessagesCount(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	// Validate token and get the user ID (UID)
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	// Optional sender ID parameter to filter unread messages by sender
	senderID := c.Query("senderId")

	// Query Firestore for all messages documents
	iter := services.FirestoreClient.Collection("messages").Documents(context.Background())

	unreadCount := 0

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var messages models.Messages
		if err := doc.DataTo(&messages); err != nil {
			log.Println("Error parsing document:", err)
			continue
		}

		for _, message := range messages.DirectMessages {
			// Check if the message is for the logged-in user
			if message.ReceiverID == uid {
				// If senderID is provided, filter by senderID
				if senderID != "" && message.SenderID != senderID {
					continue
				}

				// Check if the ReadStatus map contains the user's ID and if the value is false
				if read, exists := message.ReadStatus[uid]; !exists || !read {
					unreadCount++
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"unreadCount": unreadCount,
	})
}

// MarkMessagesAsRead updates the read status of direct messages for the logged-in user
func MarkMessagesAsRead(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in.",
		})
	}

	// Validate token and get the user ID (UID)
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again.",
		})
	}

	var payload struct {
		SenderID string `json:"senderId"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload.",
		})
	}

	if payload.SenderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Sender ID is required.",
		})
	}

	// Create a sorted list of sender and receiver IDs for consistent channel ID
	ids := []string{uid, payload.SenderID}
	sort.Strings(ids)
	channelID := strings.Join(ids, "-")

	docRef := services.FirestoreClient.Collection("messages").Doc(channelID)

	// Fetch the document from Firestore
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil || !docSnapshot.Exists() {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No messages found for the provided sender and receiver IDs.",
		})
	}

	// Convert Firestore data to the Messages struct
	var messages models.Messages
	if err := docSnapshot.DataTo(&messages); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse messages.",
		})
	}

	// Update the ReadStatus for all messages where the receiver is the logged-in user
	updated := false
	for i, message := range messages.DirectMessages {
		if message.ReceiverID == uid {
			if read, exists := message.ReadStatus[uid]; !exists || !read {
				messages.DirectMessages[i].ReadStatus[uid] = true
				updated = true
			}
		}
	}

	if !updated {
		return c.JSON(fiber.Map{
			"success": "No unread messages to mark as read.",
		})
	}

	// Update the document in Firestore
	_, err = docRef.Set(context.Background(), mappers.MapMessagesGoToFirestore(messages))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update read status.",
		})
	}

	// Trigger a Pusher event to notify the frontend to update the unread count
	notificationChannel := "user-notifications-" + uid
	err = services.PusherClient.Trigger(
		notificationChannel,
		"message-read",
		map[string]interface{}{
			"timestamp": time.Now(),
		},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger notification event.",
		})
	}

	return c.JSON(fiber.Map{
		"success": "Messages marked as read successfully.",
	})
}
