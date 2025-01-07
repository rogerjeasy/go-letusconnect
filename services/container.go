package services

import "cloud.google.com/go/firestore"

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
	// Add other services as needed
}

func NewServiceContainer(firestoreClient *firestore.Client, userSerrvice *UserService) *ServiceContainer {
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
		ChatGPTService:      NewChatGPTService(firestoreClient),
		// UserConnectionService: NewUserConnectionService(firestoreClient, userSerrvice),
		// Initialize other services
	}
}
