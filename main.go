package main

import (
	"fmt"
	"log"
	"os/signal"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/routes"
	"github.com/rogerjeasy/go-letusconnect/services"

	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/joho/godotenv"
)

func main() {
	config.LoadConfig()

	// Initialize Firebase
	if err := services.InitializeFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize Pusher
	services.InitializePusher()

	// Initialize Cloudinary
	services.InitCloudinary()

	// notificationService := services.NewNotificationService(services.FirestoreClient)
	// connectionService := services.NewUserConnectionService(services.FirestoreClient)
	userService := services.NewUserService(services.FirestoreClient)

	serviceContainer := &services.ServiceContainer{
		UserService:         userService,
		ConnectionService:   services.NewUserConnectionService(services.FirestoreClient, userService), // Pass userService here
		NotificationService: services.NewNotificationService(services.FirestoreClient),
		AuthService:         services.NewAuthService(services.FirestoreClient),
		FAQService:          services.NewFAQService(services.FirestoreClient),
		ProjectCoreService:  services.NewProjectCoreService(services.FirestoreClient),
		ProjectService:      services.NewProjectService(services.FirestoreClient, userService),
		MessageService:      services.NewMessageService(services.FirestoreClient),
		GroupChatService:    services.NewGroupChatService(services.FirestoreClient),
		NewsletterService:   services.NewNewsletterService(services.FirestoreClient),
		AddressService:      services.NewAddressService(services.FirestoreClient),
		ContactUsService:    services.NewContactUsService(services.FirestoreClient),
	}

	// Initialize ServiceContainer with all services
	// serviceContainer := services.NewServiceContainer(services.FirestoreClient, userService)

	app := fiber.New()

	// Improved CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://letusconnect.vercel.app, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Authorization",
		MaxAge:           86400,
		AllowOriginsFunc: func(origin string) bool {
			return origin == "https://letusconnect.vercel.app" ||
				origin == "http://localhost:3000"
		},
	}))

	// routes.SetupAllRoutes(app, serviceContainer)

	if err := routes.SetupAllRoutes(app, serviceContainer); err != nil {
		log.Fatalf("Failed to setup routes: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Add graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server with error handling
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
