package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTSecret is the secret key used to sign the JWT tokens.
// Ensure this key is kept private and retrieved securely (e.g., via environment variables).
var JWTSecret = []byte("your-secure-jwt-secret")

// AuthMiddleware validates JWT tokens and adds user claims to the context.
func AuthMiddleware(c *fiber.Ctx) error {
	// Extract the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// Check if the Authorization header is in the correct format
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization header format",
		})
	}

	// Extract the token from the Authorization header
	tokenString := parts[1]

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
		}
		return JWTSecret, nil
	})

	// Handle token validation errors
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Extract claims and store them in the request context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Optionally, you can perform additional validation on claims
		c.Locals("user", claims)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Proceed to the next handler
	return c.Next()
}
