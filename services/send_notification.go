package services

import (
	"context"
	"fmt"
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

func SendNewUserNotification(ctx context.Context, user *models.User) error {
	// Create a timeout context to ensure the function doesn't hang
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get all users from Firestore
	iter := FirestoreClient.Collection("users").Documents(ctx)
	defer iter.Stop()

	targetedUsersIDsList := []string{}
	readStatus := make(map[string]bool)

	// Iterate through all users
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to fetch users: %v", err)
		}

		data := doc.Data()
		targetUID, ok := data["uid"].(string)
		if !ok {
			continue // Skip invalid users without interrupting
		}

		// Don't include the new user in the targeted users list
		if targetUID != user.UID {
			targetedUsersIDsList = append(targetedUsersIDsList, targetUID)
			readStatus[targetUID] = false
		}
	}

	// Create notification service with error handling
	notificationService := NewNotificationService(FirestoreClient)
	if notificationService == nil {
		return fmt.Errorf("failed to create notification service")
	}

	notification := models.Notification{
		UserID:          user.UID,
		ActorID:         user.UID,
		ActorName:       user.Username,
		ActorType:       "user",
		Type:            models.NotificationType("new_user"),
		Title:           user.Username + " has joined the platform",
		Content:         user.Username + " has created an account on the platform. Say hello!",
		Category:        "new_user",
		Priority:        "normal",
		Status:          "unread",
		ReadStatus:      readStatus,
		IsImportant:     true,
		TargetedUsers:   targetedUsersIDsList,
		DeliveryChannel: "push",
	}

	// Save the notification with error handling
	_, err := notificationService.CreateNotification(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %v", err)
	}

	return nil
}

func SendNewGroupMessageNotification(ctx context.Context, senderID, senderName, content, groupChatID string, participants []models.Participant) error {
	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	readStatus := make(map[string]bool)
	targetedUsersIDsList := []string{}
	for _, participant := range participants {
		if participant.UserID == "" {
			return fmt.Errorf("participant UserID cannot be empty")
		}
		if participant.UserID == senderID {
			readStatus[participant.UserID] = true
		} else {
			readStatus[participant.UserID] = false
			targetedUsersIDsList = append(targetedUsersIDsList, participant.UserID)
		}
	}

	// Create notification service with error handling
	notificationService := NewNotificationService(FirestoreClient)
	if notificationService == nil {
		return fmt.Errorf("failed to create notification service")
	}

	notification := models.Notification{
		UserID:          senderID,
		ActorID:         senderID,
		ActorName:       senderName,
		ActorType:       "user",
		Type:            models.NotificationType("message"),
		Title:           "New message from " + senderName,
		Content:         content,
		Category:        "message",
		Priority:        "normal",
		Status:          "unread",
		ReadStatus:      readStatus,
		IsImportant:     true,
		GroupID:         groupChatID,
		TargetedUsers:   targetedUsersIDsList,
		DeliveryChannel: "push",
	}

	// Save the notification with error handling
	_, err := notificationService.CreateNotification(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %v", err)
	}

	return nil
}

func SendConnectionRequestNotification(ctx context.Context, fromUID, fromUsername, message, toUID string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	readStatus := map[string]bool{toUID: false}

	notificationService := NewNotificationService(FirestoreClient)
	if notificationService == nil {
		return fmt.Errorf("failed to create notification service")
	}

	if message == "" {
		message = fromUsername + " would like to connect with you"
	}

	notification := models.Notification{
		UserID:          fromUID,
		ActorID:         fromUID,
		ActorName:       fromUsername,
		ActorType:       "user",
		Type:            models.NotificationType("connection_request"),
		Title:           fromUsername + " sent you a connection request",
		Content:         message,
		Category:        "connection",
		Priority:        "normal",
		Status:          "unread",
		ReadStatus:      readStatus,
		IsImportant:     true,
		TargetedUsers:   []string{toUID},
		DeliveryChannel: "push",
	}

	_, err := notificationService.CreateNotification(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %v", err)
	}
	return nil
}

func SendConnectionAcceptedNotification(ctx context.Context, fromUID, fromUsername, toUsername, toUID string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	readStatus := map[string]bool{toUID: false}

	notificationService := NewNotificationService(FirestoreClient)
	if notificationService == nil {
		return fmt.Errorf("failed to create notification service")
	}

	notification := models.Notification{
		UserID:          fromUID,
		ActorID:         fromUID,
		ActorName:       fromUsername,
		ActorType:       "user",
		Type:            models.NotificationType("connection_accepted"),
		Title:           toUsername + " accepted your connection request",
		Content:         toUsername + " has accepted your connection request",
		Category:        "connection",
		Priority:        "normal",
		Status:          "unread",
		ReadStatus:      readStatus,
		IsImportant:     true,
		TargetedUsers:   []string{toUID},
		DeliveryChannel: "push",
	}

	_, err := notificationService.CreateNotification(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %v", err)
	}
	return nil
}
