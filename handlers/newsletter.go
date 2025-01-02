package handlers

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

type NewsletterHandler struct {
	newsletterService *services.NewsletterService
}

func NewNewsletterHandler(newsletterService *services.NewsletterService) *NewsletterHandler {
	return &NewsletterHandler{
		newsletterService: newsletterService,
	}
}

// SubscribeNewsletter handles newsletter subscriptions
func (s *NewsletterHandler) SubscribeNewsletter(c *fiber.Ctx) error {
	// Parse the request body for the email
	var requestData struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	email := strings.TrimSpace(requestData.Email)

	// Validate email format
	if _, err := mail.ParseAddress(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	ctx := context.Background()
	newsletterCollection := services.FirestoreClient.Collection("newsletters")

	// Check if the email already exists in the newsletter collection
	emailQuery := newsletterCollection.Where("email", "==", email).Documents(ctx)
	defer emailQuery.Stop()

	if _, err := emailQuery.Next(); err != iterator.Done {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "This email is already subscribed to the newsletter",
		})
	}

	// Add the email to the newsletter collection
	_, _, err := newsletterCollection.Add(ctx, map[string]interface{}{
		"email":        email,
		"subscribedAt": time.Now().Format(time.RFC3339),
	})

	if err != nil {
		log.Printf("Error subscribing email to newsletter: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to subscribe to the newsletter",
		})
	}

	// Send subscription email with an unsubscribe link
	unsubscribeLink := fmt.Sprintf("https://letusconnect.vercel.app/unsubscribe?email=%s", email)
	if err := SendNewsletterSubscriptionEmail(email, "Subscriber", unsubscribeLink); err != nil {
		log.Printf("Error sending subscription email: %v", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "You have successfully subscribed to the newsletter!",
	})
}

// UnsubscribeNewsletter handles removing a user's email from the newsletter subscriptions
func (s *NewsletterHandler) UnsubscribeNewsletter(c *fiber.Ctx) error {
	// Parse the request body for the email
	var requestData struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	email := strings.TrimSpace(requestData.Email)

	// Validate email format
	if _, err := mail.ParseAddress(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email format",
		})
	}

	ctx := context.Background()
	newsletterCollection := services.FirestoreClient.Collection("newsletters")

	// Find and delete the document with the specified email
	emailQuery := newsletterCollection.Where("email", "==", email).Documents(ctx)
	defer emailQuery.Stop()

	doc, err := emailQuery.Next()
	if err == iterator.Done {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Email not found in the newsletter subscriptions",
		})
	}

	_, err = newsletterCollection.Doc(doc.Ref.ID).Delete(ctx)
	if err != nil {
		log.Printf("Error unsubscribing email from newsletter: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unsubscribe from the newsletter",
		})
	}

	// Send unsubscription email
	if err := SendNewsletterUnsubscriptionEmail(email, "Subscriber"); err != nil {
		log.Printf("Error sending unsubscription email: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "You have successfully unsubscribed from the newsletter.",
	})
}

// GetAllSubscribers retrieves a list of all subscribed users
func (s *NewsletterHandler) GetAllSubscribers(c *fiber.Ctx) error {
	ctx := context.Background()
	newsletterCollection := services.FirestoreClient.Collection("newsletters")

	// Get all documents from the newsletter collection
	iter := newsletterCollection.Documents(ctx)
	defer iter.Stop()

	var subscribers []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error fetching subscribers: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch subscribers",
			})
		}

		subscribers = append(subscribers, doc.Data())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"subscribers": subscribers,
	})
}

// GetTotalSubscribers returns the total number of subscribed users
func (s *NewsletterHandler) GetTotalSubscribers(c *fiber.Ctx) error {
	ctx := context.Background()
	newsletterCollection := services.FirestoreClient.Collection("newsletters")

	// Count all documents in the newsletter collection
	iter := newsletterCollection.Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error counting subscribers: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to count subscribers",
			})
		}
		count++
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"totalSubscribers": count,
	})
}

// SendNewsletterSubscriptionEmail sends a welcome email to users when they subscribe to the newsletter
func SendNewsletterSubscriptionEmail(toEmail, userName, unsubscribeLink string) error {
	subject := "🎉 Thank You for Subscribing to the LetUsConnect Newsletter!"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #4A90E2;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .highlight {
      color: #E67E22;
      font-weight: bold;
    }
    .button {
      display: inline-block;
      background-color: #4A90E2;
      color: #ffffff;
      padding: 12px 24px;
      margin: 20px 0;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
    }
    .unsubscribe {
      color: #E74C3C;
      font-size: 14px;
      margin-top: 20px;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Welcome to the LetUsConnect Newsletter, %s! 📬</h2>
    <p>
      Thank you for subscribing to our newsletter! 🎉<br>
      You’re now part of our growing community of innovators, collaborators, and learners.
    </p>
    <p>
      Stay up to date with the <span class="highlight">latest developments</span>, <span class="highlight">exclusive opportunities</span>, and <span class="highlight">platform updates</span> delivered straight to your inbox.
    </p>
    <p>
      We promise to keep you informed and inspired! 💡
    </p>
    <a href="https://letusconnect.vercel.app" class="button">Explore LetUsConnect</a>

    <p class="unsubscribe">
      If you wish to unsubscribe at any time, click <a href="%s" style="color: #E74C3C;">here</a>.
    </p>
    <p>
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName, unsubscribeLink)

	return sendEmail(toEmail, subject, body)
}

// SendNewsletterUnsubscriptionEmail sends a farewell email to users when they unsubscribe from the newsletter
func SendNewsletterUnsubscriptionEmail(toEmail, userName string) error {
	subject := "Sad to See You Go! 😢😔"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <style>
    body {
      font-family: Arial, sans-serif;
      color: #333333;
      background-color: #f9f9f9;
      padding: 20px;
      text-align: center;
    }
    .container {
      max-width: 600px;
      margin: 0 auto;
      background: #ffffff;
      padding: 30px;
      border-radius: 10px;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }
    h2 {
      color: #E74C3C;
    }
    p {
      font-size: 16px;
      line-height: 1.6;
    }
    .button {
      display: inline-block;
      background-color: #4A90E2;
      color: #ffffff;
      padding: 12px 24px;
      margin: 20px 0;
      text-decoration: none;
      border-radius: 5px;
      font-weight: bold;
    }
    .footer {
      margin-top: 30px;
      font-size: 14px;
      color: #777777;
    }
  </style>
</head>
<body>
  <div class="container">
    <h2>Goodbye for Now, %s! 👋</h2>
    <p>
      We're sorry to see you go. You've successfully unsubscribed from the LetUsConnect newsletter.
    </p>
    <p>
      We hope to welcome you back someday and keep you updated with our latest news and opportunities. If you change your mind, we're always here for you.
    </p>
    <a href="https://letusconnect.vercel.app/newsletter" class="button">Resubscribe Anytime</a>

    <p class="footer">
      Best regards, <br>
      <strong>The LetUsConnect Team</strong>
    </p>
  </div>
</body>
</html>`, userName)

	return sendEmail(toEmail, subject, body)
}
