package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/services"
	"github.com/rogerjeasy/go-letusconnect/utils"
)

// PusherAuth handles Pusher authentication for private or encrypted channels
func PusherAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required. Please log in",
		})
	}

	_, err := utils.ValidateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token. Please log in again",
		})
	}

	socketID := c.FormValue("socket_id")
	channelName := c.FormValue("channel_name")

	if socketID == "" || channelName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "socket_id and channel_name are required",
		})
	}

	// Construct the string to sign
	stringToSign := fmt.Sprintf("%s:%s", socketID, channelName)
	signature := hmac.New(sha256.New, []byte(services.PusherClient.Secret))
	_, err = signature.Write([]byte(stringToSign))
	if err != nil {
		log.Printf("Error writing to HMAC: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error while generating signature",
		})
	}

	signatureBytes := signature.Sum(nil)
	authSignature := fmt.Sprintf("%s:%x", services.PusherClient.Key, signatureBytes)

	response := map[string]string{
		"auth": authSignature,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
