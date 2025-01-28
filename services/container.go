package services

import (
	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/models"
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
	WebSocketService            *WebSocketService
	UserSchoolExperienceService *UserSchoolExperienceService
	GroupService                *GroupService
	ForumService                *ForumService
	// Add other services as needed
}

func NewServiceContainer(firestoreClient *firestore.Client, userSerrvice *UserService, wsManager *models.Manager, cloudinary *cloudinary.Cloudinary) *ServiceContainer {
	pdfService := NewPDFService(firestoreClient, config.PDFContextURL)
	uploadPdfService, _ := NewUploadPDFService(firestoreClient, config.CloudinaryURL)
	ws := NewWebSocketService(wsManager)

	return &ServiceContainer{
		UserService:                 NewUserService(firestoreClient),
		ConnectionService:           NewUserConnectionService(firestoreClient, userSerrvice),
		NotificationService:         NewNotificationService(firestoreClient),
		MessageService:              NewMessageService(firestoreClient),
		GroupChatService:            NewGroupChatService(firestoreClient),
		AuthService:                 NewAuthService(firestoreClient),
		FAQService:                  NewFAQService(firestoreClient),
		ProjectCoreService:          NewProjectCoreService(firestoreClient),
		ProjectService:              NewProjectService(firestoreClient, userSerrvice),
		AddressService:              NewAddressService(firestoreClient),
		NewsletterService:           NewNewsletterService(firestoreClient),
		ContactUsService:            NewContactUsService(firestoreClient),
		PDFService:                  pdfService,
		ChatGPTService:              NewChatGPTService(firestoreClient, pdfService),
		UploadPDFService:            uploadPdfService,
		WebSocketService:            ws,
		UserSchoolExperienceService: NewUserSchoolExperienceService(firestoreClient, userSerrvice),
		GroupService:                NewGroupService(firestoreClient, cloudinary, userSerrvice),
		ForumService:                NewForumService(firestoreClient, userSerrvice),
		// WebSocketService:    NewWebSocketService(firestoreClient),
		// UserConnectionService: NewUserConnectionService(firestoreClient, userSerrvice),
		// Initialize other services
	}
}
