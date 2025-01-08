package services

import (
	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/config"
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
	// Add other services as needed
}

func NewServiceContainer(firestoreClient *firestore.Client, userSerrvice *UserService) *ServiceContainer {
	pdfService := NewPDFService(firestoreClient, config.PDFContextURL)
	return &ServiceContainer{
		UserService:         NewUserService(firestoreClient),
		ConnectionService:   NewUserConnectionService(firestoreClient, userSerrvice),
		NotificationService: NewNotificationService(firestoreClient),
		MessageService:      NewMessageService(firestoreClient),
		GroupChatService:    NewGroupChatService(firestoreClient),
		AuthService:         NewAuthService(firestoreClient),
		FAQService:          NewFAQService(firestoreClient),
		ProjectCoreService:  NewProjectCoreService(firestoreClient),
		ProjectService:      NewProjectService(firestoreClient, userSerrvice),
		AddressService:      NewAddressService(firestoreClient),
		NewsletterService:   NewNewsletterService(firestoreClient),
		ContactUsService:    NewContactUsService(firestoreClient),
		PDFService:          pdfService,
		ChatGPTService:      NewChatGPTService(firestoreClient, pdfService),
		// UserConnectionService: NewUserConnectionService(firestoreClient, userSerrvice),
		// Initialize other services
	}
}
