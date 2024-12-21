package main

import (
	"log"

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

	app := fiber.New()

	// Enable CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://letusconnect.vercel.app",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server on the port provided by Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
