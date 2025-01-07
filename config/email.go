package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SMTPHost       string
	SMTPPort       string
	SenderEmail    string
	SenderName     string
	SenderPass     string
	FirebaseAPiKey string

	FirebaseSignInURL      string
	EnvServiceAccount      string
	JsonServiceAccountPath string

	PusherAppID   string
	PusherKey     string
	PusherSecret  string
	PusherCluster string

	CloudinaryURL string
	OpenAIKey     string
)

func LoadConfig() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}

	// Retrieve values from environment variables
	SMTPHost = os.Getenv("SMTP_HOST")
	SMTPPort = os.Getenv("SMTP_PORT")
	SenderEmail = os.Getenv("SENDER_EMAIL")
	SenderName = os.Getenv("SENDER_NAME")
	SenderPass = os.Getenv("SENDER_PASS")
	FirebaseAPiKey = os.Getenv("FIREBASE_API_KEY")
	EnvServiceAccount = os.Getenv("ENV_SERVICE_ACCOUNT_KEY")
	JsonServiceAccountPath = os.Getenv("JSON_SERVICE_ACCOUNT_PATH")

	PusherAppID = os.Getenv("PUSHER_APP_ID")
	PusherKey = os.Getenv("PUSHER_KEY")
	PusherSecret = os.Getenv("PUSHER_SECRET")
	PusherCluster = os.Getenv("PUSHER_CLUSTER")

	// Initialize FirebaseSignInURL after loading the API key
	FirebaseSignInURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=" + FirebaseAPiKey

	CloudinaryURL = os.Getenv("CLOUDINARY_URL")

	// openai
	OpenAIKey = os.Getenv("OPENAI_API_KEY")
}
