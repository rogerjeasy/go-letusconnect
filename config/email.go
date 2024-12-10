package config

import (
	"os"
)

// Email configuration constants
var (
	SMTPHost    = os.Getenv("SMTP_HOST")    // e.g., "smtp.gmail.com"
	SMTPPort    = os.Getenv("SMTP_PORT")    // e.g., "587"
	SenderEmail = os.Getenv("SENDER_EMAIL") // e.g., "your-email@gmail.com"
	SenderPass  = os.Getenv("SENDER_PASS")  // App-specific password or SMTP password
	SenderName  = "LetUsConnect Support"    // Name displayed in the email
)
