package services

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/services/sms"
)

type ServiceContainer struct {
	UserService           *UserService
	ConnectionService     *UserConnectionService
	NotificationService   *NotificationService
	MessageService        *MessageService
	GroupChatService      *GroupChatService
	AuthService           *AuthService
	FAQService            *FAQService
	ProjectCoreService    *ProjectCoreService
	ProjectService        *ProjectService
	UserConnectionService *UserConnectionService
	AddressService        *AddressService
	NewsletterService     *NewsletterService
	ContactUsService      *ContactUsService
	ChatGPTService        *ChatGPTService
	PDFService            *PDFService
	UploadPDFService      *UploadPDFService
	// WebSocketService      *WebSocketService
	WebSocketService             *WebSocketService
	UserSchoolExperienceService  *UserSchoolExperienceService
	GroupService                 *GroupService
	ForumService                 *ForumService
	TestimonialService           *TestimonialService
	GeneralNotificationService   *GeneralNotificationService
	JobService                   *JobService
	LinkedInJobsService          *LinkedInJobsService
	SchedulerNotificationService *SchedulerNotificationService
	notificationScheduler        *NotificationScheduler
	// Add other services as needed
}

func NewServiceContainer(firestoreClient FirestoreClient, userSerrvice *UserService, cloudinary *cloudinary.Cloudinary) *ServiceContainer {
	pdfService := NewPDFService(firestoreClient, config.PDFContextURL)
	uploadPdfService, _ := NewUploadPDFService(firestoreClient, config.CloudinaryURL)

	// Initialize SMS service
	smsService := sms.NewSMSService(sms.Config{
		AccountSID: config.TwilioAccountSID,
		AuthToken:  config.TwilioAuthToken,
		FromNumber: config.TwilioFromNumber,
	})

	// Initialize notification scheduler
	notificationScheduler := NewNotificationScheduler(firestoreClient, smsService)

	return &ServiceContainer{
		UserService:                  NewUserService(firestoreClient),
		ConnectionService:            NewUserConnectionService(firestoreClient, userSerrvice),
		NotificationService:          NewNotificationService(firestoreClient),
		MessageService:               NewMessageService(firestoreClient),
		GroupChatService:             NewGroupChatService(firestoreClient),
		AuthService:                  NewAuthService(firestoreClient),
		FAQService:                   NewFAQService(firestoreClient),
		ProjectCoreService:           NewProjectCoreService(firestoreClient),
		ProjectService:               NewProjectService(firestoreClient, userSerrvice),
		AddressService:               NewAddressService(firestoreClient),
		NewsletterService:            NewNewsletterService(firestoreClient),
		ContactUsService:             NewContactUsService(firestoreClient),
		PDFService:                   pdfService,
		ChatGPTService:               NewChatGPTService(firestoreClient, pdfService),
		UploadPDFService:             uploadPdfService,
		UserSchoolExperienceService:  NewUserSchoolExperienceService(firestoreClient, userSerrvice),
		GroupService:                 NewGroupService(firestoreClient, cloudinary, userSerrvice),
		ForumService:                 NewForumService(firestoreClient, userSerrvice),
		TestimonialService:           NewTestimonialService(firestoreClient, userSerrvice),
		GeneralNotificationService:   NewGeneralNotificationService(firestoreClient),
		JobService:                   NewJobService(firestoreClient),
		LinkedInJobsService:          NewLinkedInJobsService(firestoreClient),
		notificationScheduler:        notificationScheduler,
		SchedulerNotificationService: NewSchedulerNotificationService(notificationScheduler),
		// WebSocketService:    NewWebSocketService(firestoreClient),
		// UserConnectionService: NewUserConnectionService(firestoreClient, userSerrvice),
		// Initialize other services
	}
}

// StartServices initializes and starts any background services
// func (sc *ServiceContainer) StartServices(ctx context.Context) {
// 	// Start the notification scheduler
// 	if sc.notificationScheduler != nil {
// 		sc.notificationScheduler.Start(ctx)
// 	}
// }

// StopServices gracefully shuts down any running services
// func (sc *ServiceContainer) StopServices() {
// 	// Stop the notification scheduler
// 	if sc.notificationScheduler != nil {
// 		sc.notificationScheduler.Stop()
// 	}
// }
