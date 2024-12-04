package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/routes"
	"github.com/rogerjeasy/go-letusconnect/services"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}

func main() {
	// Initialize Firebase
	if err := services.InitializeFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	app := fiber.New()

	// Enable CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000", // Allow requests from your frontend
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
