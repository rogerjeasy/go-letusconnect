package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

// Secret key for JWT (use environment variables in production)
var jwtSecretKey = []byte("your_jwt_secret_key")

// GenerateJWT generates a new JWT token
func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"uid":   user.UID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

// Parse and validate the JWT token
func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		return "", errors.New("missing UID in token claims")
	}

	return uid, nil
}

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// Register creates a new user in Firebase Authentication and Firestore
func Register(c *fiber.Ctx) error {
	// Parse the request body into a map
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Map the frontend data to the User struct
	user := mappers.MapFrontendToUser(requestData)

	// Validate required fields
	if strings.TrimSpace(user.Username) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Username is required",
		})
	}
	if strings.TrimSpace(user.Email) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	if strings.TrimSpace(user.Password) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}
	if strings.TrimSpace(user.Program) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Program is required",
		})
	}

	// Validate uniqueness of username and email in Firestore
	ctx := context.Background()

	// Check for existing user with the same email
	emailQuery := services.FirestoreClient.Collection("users").Where("email", "==", user.Email).Documents(ctx)
	defer emailQuery.Stop()
	if _, err := emailQuery.Next(); err != iterator.Done {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Email already in use",
		})
	}

	// Check for existing user with the same username
	usernameQuery := services.FirestoreClient.Collection("users").Where("username", "==", user.Username).Documents(ctx)
	defer usernameQuery.Stop()
	if _, err := usernameQuery.Next(); err != iterator.Done {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Username already in use",
		})
	}

	// Create user in Firebase Authentication
	authUser, err := services.FirebaseAuth.CreateUser(ctx, (&auth.UserToCreate{}).
		Email(user.Email).
		Password(user.Password))
	if err != nil {
		log.Printf("Error creating user in Firebase Auth: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user in Firebase Authentication",
		})
	}

	// Add additional fields
	user.UID = authUser.UID
	currentTime := time.Now()
	customFormat := "Monday, Jan 2, 2006 at 3:04 PM"
	formattedTime := FormatTime(currentTime, customFormat)
	user.AccountCreatedAt = formattedTime
	user.IsActive = true
	user.IsVerified = false
	user.Role = []string{"user"}
	user.Password = ""

	// Convert the user struct to Firestore-compatible (snake_case) format
	backendUser := mappers.MapUserFrontendToBackend(&user)

	// Save user to Firestore
	_, _, err = services.FirestoreClient.Collection("users").Add(ctx, backendUser)
	if err != nil {
		log.Printf("Error saving user to Firestore: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save user",
		})
	}

	// Generate JWT token
	token, err := GenerateJWT(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set JWT token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	// Map backend user to frontend (camelCase) format
	frontendUser := mappers.MapUserBackendToFrontend(backendUser)

	// Return success response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    frontendUser,
		"token":   token,
	})
}

// Login authenticates a user using Firebase Authentication

// FirebaseSignInURL is the Firebase REST API endpoint for sign-in

// Login authenticates the user and returns a JWT token
// Login authenticates the user and returns a JWT token
func Login(c *fiber.Ctx) error {
	var user models.User

	// Parse request body into User model
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate input fields
	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Prepare request payload for Firebase Authentication REST API
	payload := map[string]string{
		"email":             user.Email,
		"password":          user.Password,
		"returnSecureToken": "true",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process authentication request",
		})
	}

	// Make the REST API request to Firebase
	resp, err := http.Post(FirebaseSignInURL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	defer resp.Body.Close()

	// Parse the response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var firebaseResponse map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &firebaseResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process Firebase response",
		})
	}

	// Retrieve user details from Firestore
	ctx := context.Background()
	userQuery := services.FirestoreClient.Collection("users").Where("email", "==", user.Email).Documents(ctx)
	defer userQuery.Stop()

	doc, err := userQuery.Next()
	if err == iterator.Done {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var dbUser models.User
	if err := doc.DataTo(&dbUser); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user data",
		})
	}

	// Generate JWT token
	token, err := GenerateJWT(&dbUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set JWT token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	// Map backend user to frontend format
	frontendUser := mappers.MapUserBackendToFrontend(mappers.MapUserFrontendToBackend(&dbUser))

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user":    frontendUser,
	})
}
