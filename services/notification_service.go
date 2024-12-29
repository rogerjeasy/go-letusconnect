package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotificationService handles operations related to notifications
type NotificationService struct {
	firestoreClient *firestore.Client
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(client *firestore.Client) *NotificationService {
	return &NotificationService{
		firestoreClient: client,
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	notification.ID = uuid.New().String()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()

	_, err := s.firestoreClient.Collection("notifications").Doc(notification.ID).Set(ctx, mappers.MapNotificationGoToFirestore(notification))
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

// UpdateNotification updates an existing notification
func (s *NotificationService) UpdateNotification(ctx context.Context, notificationID string, updates map[string]interface{}) (*models.Notification, error) {
	docRef := s.firestoreClient.Collection("notifications").Doc(notificationID)

	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		currentNotification := mappers.MapNotificationFirestoreToGo(doc.Data())

		for key, value := range updates {
			switch key {
			case "title":
				currentNotification.Title = value.(string)
			case "content":
				currentNotification.Content = value.(string)
			case "status":
				currentNotification.Status = models.NotificationStatus(value.(string))
				// Add more fields as needed
			}
		}

		currentNotification.UpdatedAt = time.Now()

		return tx.Set(docRef, mappers.MapNotificationGoToFirestore(currentNotification))
	})

	if err != nil {
		return nil, err
	}

	updatedDoc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	updatedNotification := mappers.MapNotificationFirestoreToGo(updatedDoc.Data())
	return &updatedNotification, nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID string) error {
	_, err := s.firestoreClient.Collection("notifications").Doc(notificationID).Delete(ctx)
	return err
}

// GetNotification fetches a single notification by ID
func (s *NotificationService) GetNotification(ctx context.Context, notificationID string) (*models.Notification, error) {
	doc, err := s.firestoreClient.Collection("notifications").Doc(notificationID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}

	notification := mappers.MapNotificationFirestoreToGo(doc.Data())
	return &notification, nil
}

// ListNotifications fetches notifications for a user
func (s *NotificationService) ListNotifications(ctx context.Context, userID string, limit int, lastNotificationID string) ([]models.Notification, error) {
	// query := s.firestoreClient.Collection("notifications").Where("user_id", "==", userID).OrderBy("created_at", firestore.Desc).Limit(limit)
	query := s.firestoreClient.Collection("notifications").Where("user_id", "==", userID)

	if lastNotificationID != "" {
		lastDoc, err := s.firestoreClient.Collection("notifications").Doc(lastNotificationID).Get(ctx)
		if err != nil {
			return nil, err
		}
		query = query.StartAfter(lastDoc)
	}

	iter := query.Documents(ctx)
	var notifications []models.Notification

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		notification := mappers.MapNotificationFirestoreToGo(doc.Data())
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// MarkNotificationAsRead marks a notification as read for a specific user
func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID string, userID string) error {
	// First get the current notification to access the existing read_status map
	docRef := s.firestoreClient.Collection("notifications").Doc(notificationID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get notification: %v", err)
	}

	// Get current read status map
	readStatus := make(map[string]bool)
	currentReadStatus, exists := doc.Data()["read_status"]
	if exists && currentReadStatus != nil {
		// Type assert to map[string]interface{} first
		if readStatusMap, ok := currentReadStatus.(map[string]interface{}); ok {
			// Convert each value to bool
			for k, v := range readStatusMap {
				if boolVal, ok := v.(bool); ok {
					readStatus[k] = boolVal
				}
			}
		}
	}

	// Update read status for the specific user
	readStatus[userID] = true

	// Update the document
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "read_status", Value: readStatus},
		{Path: "updated_at", Value: time.Now()},
		{Path: "read_at", Value: time.Now()},
	})
	if err != nil {
		return fmt.Errorf("failed to update notification: %v", err)
	}

	return nil
}

// ListTargetedNotifications fetches notifications where the user is in the targeted_users list
func (s *NotificationService) ListTargetedNotifications(ctx context.Context, userID string, limit int, lastNotificationID string) ([]models.Notification, error) {
	// Create a query for notifications where userID is in targeted_users array
	query := s.firestoreClient.Collection("notifications").Where("targeted_users", "array-contains", userID)

	// Add pagination if lastNotificationID is provided
	if lastNotificationID != "" {
		lastDoc, err := s.firestoreClient.Collection("notifications").Doc(lastNotificationID).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get last notification: %v", err)
		}
		query = query.StartAfter(lastDoc)
	}

	// Add ordering and limit
	query = query.OrderBy("created_at", firestore.Desc).OrderBy("__name__", firestore.Desc)
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute query with error handling for missing index
	iter := query.Documents(ctx)
	defer iter.Stop()

	var notifications []models.Notification
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// Check if the error is due to missing index
			if status.Code(err) == codes.FailedPrecondition {
				// Fallback to unordered query if index is not ready
				return s.listTargetedNotificationsWithoutOrdering(ctx, userID, limit)
			}
			return nil, fmt.Errorf("error iterating notifications: %v", err)
		}

		notification := mappers.MapNotificationFirestoreToGo(doc.Data())

		// Initialize read status map if nil
		if notification.ReadStatus == nil {
			notification.ReadStatus = make(map[string]bool)
		}
		if _, exists := notification.ReadStatus[userID]; !exists {
			notification.ReadStatus[userID] = false
		}

		// Initialize archived status map if nil
		if notification.IsArchived == nil {
			notification.IsArchived = make(map[string]bool)
		}
		if _, exists := notification.IsArchived[userID]; !exists {
			notification.IsArchived[userID] = false
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// Fallback function for when index is not ready
func (s *NotificationService) listTargetedNotificationsWithoutOrdering(ctx context.Context, userID string, limit int) ([]models.Notification, error) {
	// Simple query without ordering
	query := s.firestoreClient.Collection("notifications").Where("targeted_users", "array-contains", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var notifications []models.Notification
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error in fallback query: %v", err)
		}

		notification := mappers.MapNotificationFirestoreToGo(doc.Data())

		// Initialize maps if nil
		if notification.ReadStatus == nil {
			notification.ReadStatus = make(map[string]bool)
		}
		if _, exists := notification.ReadStatus[userID]; !exists {
			notification.ReadStatus[userID] = false
		}

		if notification.IsArchived == nil {
			notification.IsArchived = make(map[string]bool)
		}
		if _, exists := notification.IsArchived[userID]; !exists {
			notification.IsArchived[userID] = false
		}

		notifications = append(notifications, notification)
	}

	// Sort notifications by created_at manually since we can't use Firestore ordering
	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].CreatedAt.After(notifications[j].CreatedAt)
	})

	return notifications, nil
}

// CountUnreadNotifications counts notifications where ReadStatus[userID] is false
func (s *NotificationService) CountUnreadNotifications(ctx context.Context, userID string) (int64, error) {
	query := s.firestoreClient.Collection("notifications").Where("targeted_users", "array-contains", userID)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var unreadCount int64 = 0

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("error counting unread notifications: %v", err)
		}

		notification := mappers.MapNotificationFirestoreToGo(doc.Data())

		// Simply check if ReadStatus[userID] is false
		if !notification.ReadStatus[userID] {
			unreadCount++
		}
	}

	return unreadCount, nil
}

func (s *NotificationService) GetNotificationStats(ctx context.Context, userID string) (models.NotificationStats, error) {
	query := s.firestoreClient.Collection("notifications").Where("targeted_users", "array-contains", userID)

	iter := query.Documents(ctx)
	defer iter.Stop()

	stats := models.NotificationStats{
		PriorityStats: make(map[string]int64),
		TypeStats:     make(map[string]int64),
	}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return stats, fmt.Errorf("error getting notification stats: %v", err)
		}

		notification := mappers.MapNotificationFirestoreToGo(doc.Data())
		stats.TotalCount++

		// Simply check ReadStatus[userID]
		if notification.ReadStatus[userID] {
			stats.ReadCount++
		} else {
			stats.UnreadCount++
		}

		if notification.IsArchived[userID] {
			stats.ArchivedCount++
		}

		stats.PriorityStats[string(notification.Priority)]++
		stats.TypeStats[string(notification.Type)]++
	}

	return stats, nil
}
