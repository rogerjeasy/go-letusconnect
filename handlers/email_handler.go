package handlers

import (
	"fmt"
	"net/smtp"

	"github.com/rogerjeasy/go-letusconnect/config"
)

// SendAutomaticEmail sends a thank-you email to the sender
func SendAutomaticEmail(toEmail, senderName string) error {
	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPass, config.SMTPHost)

	// Email subject and body
	subject := "Thank You for Contacting Us!"
	body := fmt.Sprintf(`Dear %s,

Thank you for reaching out to LetUsConnect Support. We have received your message and will get back to you as soon as possible.

Best regards,  
The LetUsConnect Team`, senderName)

	// Email message format
	msg := fmt.Sprintf("From: %s <%s>\r\n", config.SenderName, config.SenderEmail) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort),
		auth,
		config.SenderEmail,
		[]string{toEmail},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
