package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services/sms"
	"google.golang.org/api/iterator"
)

const NOTIFICATIONS_COLLECTION = "notifications"

type NotificationScheduler struct {
	firestoreClient FirestoreClient
	smsService      *sms.SMSService
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

func NewNotificationScheduler(
	client FirestoreClient,
	smsService *sms.SMSService,
) *NotificationScheduler {
	return &NotificationScheduler{
		firestoreClient: client,
		smsService:      smsService,
		stopChan:        make(chan struct{}),
	}
}

func (s *NotificationScheduler) Start(ctx context.Context) {
	s.wg.Add(1)
	go s.run(ctx)
}

func (s *NotificationScheduler) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}

func (s *NotificationScheduler) run(ctx context.Context) {
	defer s.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.processNotifications(ctx)
		}
	}
}

func (s *NotificationScheduler) processNotifications(ctx context.Context) {
	now := time.Now()

	// Query for pending notifications that are due
	query := s.firestoreClient.Collection(NOTIFICATIONS_COLLECTION).
		Where("status", "==", models.NotificationStatusPending).
		Where("scheduledAt", "<=", now)

	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating notifications: %v", err)
			continue
		}

		var notification models.Notification
		if err := doc.DataTo(&notification); err != nil {
			log.Printf("Error converting document to notification: %v", err)
			continue
		}

		err = s.sendNotification(&notification)
		if err != nil {
			log.Printf("Error sending notification %s: %v", notification.ID, err)
			notification.Status = models.NotificationStatusFailed
		} else {
			notification.Status = models.NotificationStatusSent
			sentTime := time.Now()
			notification.SentAt = &sentTime
		}

		// Update notification status in Firestore
		_, err = doc.Ref.Set(ctx, notification, firestore.MergeAll)
		if err != nil {
			log.Printf("Error updating notification status: %v", err)
		}
	}
}

func (s *NotificationScheduler) sendNotification(notification *models.Notification) error {
	switch notification.Type {
	case models.NotificationTypeSMS:
		return s.smsService.SendSMS(
			notification.Recipient,
			notification.Content,
		)
	case models.NotificationTypeEmail:
		return fmt.Errorf("email notifications not implemented")
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}

func (s *NotificationScheduler) ScheduleNotification(ctx context.Context, req *models.NotificationRequest) error {
	// Validate notification type
	if req.Type != models.NotificationTypeSMS && req.Type != models.NotificationTypeEmail {
		return fmt.Errorf("invalid notification type: %s", req.Type)
	}

	notification := &models.Notification{
		ID:              uuid.New().String(),
		UserID:          req.UserID,
		Type:            req.Type,
		Status:          models.NotificationStatusPending,
		Title:           req.Subject,
		Content:         req.Content,
		DeliveryChannel: string(req.Type),
		ScheduledAt:     req.ScheduledAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Priority:        models.NotificationPriorityNormal,
	}

	_, err := s.firestoreClient.Collection(NOTIFICATIONS_COLLECTION).Doc(notification.ID).Set(ctx, notification)
	return err
}

func (s *NotificationScheduler) CancelNotification(ctx context.Context, notificationID string) error {
	_, err := s.firestoreClient.Collection(NOTIFICATIONS_COLLECTION).Doc(notificationID).Update(ctx, []firestore.Update{
		{
			Path:  "status",
			Value: models.NotificationStatusCancelled,
		},
		{
			Path:  "updatedAt",
			Value: time.Now(),
		},
	})
	return err
}

func (s *NotificationScheduler) GetNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	var notifications []models.Notification

	query := s.firestoreClient.Collection(NOTIFICATIONS_COLLECTION).Where("userId", "==", userID)
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var notification models.Notification
		if err := doc.DataTo(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

// GetPendingNotifications returns all pending notifications
func (s *NotificationScheduler) GetPendingNotifications(ctx context.Context) ([]models.Notification, error) {
	var notifications []models.Notification

	query := s.firestoreClient.Collection(NOTIFICATIONS_COLLECTION).Where("status", "==", models.NotificationStatusPending)
	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var notification models.Notification
		if err := doc.DataTo(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
