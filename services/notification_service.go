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

// MarkNotificationAsRead marks a notification as read
func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID string) error {
	_, err := s.firestoreClient.Collection("notifications").Doc(notificationID).Update(ctx, []firestore.Update{
		{Path: "status", Value: string(models.NotificationStatusRead)},
		{Path: "updated_at", Value: time.Now()},
		{Path: "read_at", Value: time.Now()},
	})
	return err
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
