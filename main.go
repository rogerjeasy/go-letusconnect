package main

import (
	"fmt"
	"log"
	"os/signal"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/middleware"
	"github.com/rogerjeasy/go-letusconnect/routes"
	"github.com/rogerjeasy/go-letusconnect/services"
	// "github.com/joho/godotenv"
)

func main() {

	// Create a context that we'll use to manage service lifecycles
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	config.LoadConfig()

	// Initialize Firebase
	if err := services.InitializeFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize Pusher
	services.InitializePusher()

	userService := services.NewUserService(services.Firestore)
	cloudinary := services.InitCloudinary()

	serviceContainer := services.NewServiceContainer(services.Firestore, userService, cloudinary)

	// Start background services
	// serviceContainer.StartServices(ctx)
	// defer serviceContainer.StopServices()

	authService := services.NewAuthService(services.Firestore)
	authHandler := handlers.NewAuthHandler(authService, serviceContainer)

	app := fiber.New()

	app.Post("/auth/linkedin", authHandler.LinkedInCallback)

	websocketService := services.NewGorillaWebSocketService()
	websocketService.HandleWebSocket(app)

	// Improved CORS configuration
	app.Use(middleware.ConfigureCORS())

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
		// Cleanup PDFService
		if serviceContainer.PDFService != nil {
			serviceContainer.PDFService.Stop()
		}
		_ = app.Shutdown()
	}()

	// Start server with error handling
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
