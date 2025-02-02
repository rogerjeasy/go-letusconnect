package handlers

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

type ContactUsHandler struct {
	contactUsService *services.ContactUsService
}

func NewContactUsHandler(contactUsService *services.ContactUsService) *ContactUsHandler {
	return &ContactUsHandler{
		contactUsService: contactUsService,
	}
}

// CreateContact handles creating a new contact form submission
func (h *ContactUsHandler) CreateContact(c *fiber.Ctx) error {
	var requestData map[string]interface{}

	// Parse the request body
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()

	// Map the request data to the ContactUs struct
	newContact := mappers.FrontendToContactUs(requestData)

	// Save to Firestore
	docRef, _, err := services.Firestore.Collection("contact_us").Add(ctx, mappers.ContactUsToFirestore(newContact))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create contact",
		})
	}

	newContact.ID = docRef.ID

	// Send automatic thank-you email
	if err := SendAutomaticEmail(newContact.Email, newContact.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Contact created, but failed to send confirmation email",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Thank you for reaching out to us! Your message has been received, and our team will get back to you as soon as possible. We appreciate your patience and look forward to assisting you.",
		"data":    newContact,
	})
}

// GetContact retrieves a contact form submission by ID
func (h *ContactUsHandler) GetContactByID(c *fiber.Ctx) error {
	contactID := c.Params("id")
	if contactID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Contact ID is required",
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate the token and extract the UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	ctx := context.Background()

	// Fetch the document by ID
	doc, err := services.Firestore.Collection("contact_us").Doc(contactID).Get(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Contact not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch contact",
		})
	}

	// Map the Firestore document data to the ContactUs struct
	contact := mappers.FirestoreToContactUs(doc.Data())

	return c.Status(fiber.StatusOK).JSON(contact)
}

// GetAllContacts retrieves all contact form submissions
func (h *ContactUsHandler) GetAllContacts(c *fiber.Ctx) error {

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate the token and extract the UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	ctx := context.Background()

	// Query Firestore to get all contacts
	iter := services.Firestore.Collection("contact_us").Documents(ctx)
	defer iter.Stop()

	var contacts []models.ContactUs

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch contacts",
			})
		}

		contact := mappers.FirestoreToContactUs(doc.Data())
		contacts = append(contacts, *contact)
	}

	return c.Status(fiber.StatusOK).JSON(contacts)
}

// UpdateContactStatus updates the status of a contact form submission (e.g., to "read" or "replied")
func (h *ContactUsHandler) UpdateContactStatus(c *fiber.Ctx) error {
	contactID := c.Params("id")
	if contactID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Contact ID is required",
		})
	}

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate the token and extract the UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var requestData struct {
		Status       string `json:"status"`
		RepliedBy    string `json:"repliedBy,omitempty"`
		ReplyMessage string `json:"replyMessage,omitempty"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()

	updates := []firestore.Update{
		{Path: "status", Value: requestData.Status},
		{Path: "updated_at", Value: time.Now()},
	}

	if requestData.RepliedBy != "" {
		updates = append(updates, firestore.Update{Path: "replied_by", Value: requestData.RepliedBy})
	}

	if requestData.ReplyMessage != "" {
		updates = append(updates, firestore.Update{Path: "reply_message", Value: requestData.ReplyMessage})
	}

	// Update Firestore document
	_, err_new := services.Firestore.Collection("contacts").Doc(contactID).Update(ctx, updates)
	if err_new != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update contact status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Contact status updated successfully",
	})
}
