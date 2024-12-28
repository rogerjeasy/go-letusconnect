package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	// "os"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

var jwtSecretKey = []byte("your_jwt_secret_key")

func generateRandomAvatar() string {

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	uniqueID := fmt.Sprintf("%x", rng.Int63())

	return fmt.Sprintf("https://picsum.photos/seed/%s/150/150?nature", uniqueID)
}

func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"uid":   user.UID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func validateToken(tokenString string) (string, error) {
	if strings.TrimSpace(tokenString) == "" {
		log.Printf("Empty token received")
		return "", errors.New("token cannot be empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Invalid signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		log.Printf("Token is invalid")
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Failed to extract token claims")
		return "", errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			log.Printf("Token expired at %v", time.Unix(int64(exp), 0))
			return "", errors.New("token has expired")
		}
	} else {
		log.Printf("Token missing expiration claim")
		return "", errors.New("invalid token: missing expiration")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		log.Printf("Missing or invalid UID in token claims")
		return "", errors.New("missing or invalid UID in token")
	}

	if strings.TrimSpace(uid) == "" {
		log.Printf("Empty UID in token claims")
		return "", errors.New("empty UID in token")
	}

	return uid, nil
}

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

func Register(c *fiber.Ctx) error {
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

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

	randomAvatarURL := generateRandomAvatar()

	// Upload the generated avatar to Cloudinary
	cld := services.CloudinaryClient
	uploadResult, err := cld.Upload.Upload(ctx, randomAvatarURL, uploader.UploadParams{
		PublicID: fmt.Sprintf("users/%s/avatar", user.Username),
		Folder:   "users/avatars",
	})
	if err != nil {
		log.Printf("Error uploading avatar to Cloudinary: %v", err)
		user.ProfilePicture = randomAvatarURL
	} else {
		user.ProfilePicture = uploadResult.SecureURL
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
	user.AccountCreatedAt = FormatTime(currentTime, customFormat)
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

	// Send welcome email
	err = SendWelcomeEmail(user.Email, user.Username, "")
	if err != nil {
		log.Printf("Error sending welcome email: %v", err)
		// Don't fail the registration process if email sending fails
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

	// After successfully creating the user
	go func() {
		if err := services.SendNewUserNotification(context.Background(), &user); err != nil {
			log.Printf("Failed to send new user notification: %v", err)
		}
	}()

	// Return success response
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "You have successfully created an account",
		"user":    frontendUser,
		"token":   token,
	})
}

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
	resp, err := http.Post(config.FirebaseSignInURL, "application/json", bytes.NewReader(payloadBytes))
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

	var dbUser map[string]interface{}
	if err := doc.DataTo(&dbUser); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user data",
		})
	}

	// Convert Firestore data to the `models.User` struct
	backendUser := mappers.MapBackendToUser(dbUser)

	// Generate JWT token
	token, err := GenerateJWT(&backendUser)
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
	frontendUser := mappers.MapUserToFrontend(&backendUser)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "You have successfully logged in to your account",
		"token":   token,
		"user":    frontendUser,
	})
}

func Logout(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

// FetchPlatformLogoURL retrieves the platform's logo URL from Firestore
func FetchPlatformLogoURL() (string, error) {
	ctx := context.Background()
	logoDoc, err := services.FirestoreClient.Collection("config").Doc("platform").Get(ctx)
	if err != nil {
		log.Printf("Error fetching platform logo: %v", err)
		return "", err
	}

	// Retrieve the logo URL from the document
	logoURL, ok := logoDoc.Data()["logo_url"].(string)
	if !ok {
		log.Printf("No logo URL found in Firestore")
		return "", nil
	}

	return logoURL, nil
}
