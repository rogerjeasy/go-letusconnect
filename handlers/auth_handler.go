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

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

type AuthHandler struct {
	authService      *services.AuthService
	containerService *services.ServiceContainer
}

func NewAuthHandler(authService *services.AuthService, containerService *services.ServiceContainer) *AuthHandler {
	return &AuthHandler{
		authService:      authService,
		containerService: containerService,
	}
}

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

func (a *AuthHandler) Register(c *fiber.Ctx) error {
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	providerType := models.AuthProvider(requestData["provider"].(string))
	if providerType == "" {
		providerType = models.EmailPassword
	}

	providerData, err := extractProviderData(providerType, requestData)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate common required fields
	if err := validateCommonFields(providerData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ctx := context.Background()

	// Check for existing user
	if err := checkExistingUser(ctx, providerData); err != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var profilePictureURL string
	if providerData.PhotoURL != "" {
		profilePictureURL = providerData.PhotoURL
	} else {
		profilePictureURL = generateRandomAvatar()
	}

	// Upload profile picture to Cloudinary
	uploadedURL, err := uploadProfilePicture(ctx, profilePictureURL, providerData.Username)
	if err != nil {
		log.Printf("Error uploading to Cloudinary: %v", err)
		uploadedURL = profilePictureURL
	}

	// Create or get Firebase user
	authUser, err := createFirebaseUser(ctx, providerData)
	if err != nil {
		log.Printf("Error with Firebase auth: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate user",
		})
	}

	// Create user model
	currentTime := time.Now()
	customFormat := "Monday, Jan 2, 2006 at 3:04 PM"

	user := models.User{
		UID:              authUser.UID,
		Username:         providerData.Username,
		FirstName:        providerData.FirstName,
		LastName:         providerData.LastName,
		Email:            providerData.Email,
		ProfilePicture:   uploadedURL,
		Program:          providerData.Program,
		AccountCreatedAt: FormatTime(currentTime, customFormat),
		IsActive:         true,
		IsVerified:       providerType != models.EmailPassword,
		Role:             []string{"user"},
		IsOnline:         true,
		Bio:              "",
		PhoneNumber:      "",
		GraduationYear:   0,
		Interests:        []string{},
		Skills:           []string{},
		Languages:        []string{},
		Projects:         []string{},
		Certifications:   []string{},
		IsPrivate:        false,
	}

	// Convert to backend format and save to Firestore
	if err := a.authService.CreateUser(ctx, &user); err != nil {
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

	// Set JWT cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	// Send welcome email asynchronously
	go func() {
		if err := SendWelcomeEmail(user.Email, user.Username, string(providerType)); err != nil {
			log.Printf("Error sending welcome email: %v", err)
		}
	}()

	// Send new user notification asynchronously
	go func() {
		if err := a.containerService.GeneralNotificationService.SendNewUserNotification(context.Background(), &user); err != nil {
			log.Printf("Failed to send new user notification: %v", err)
		}
	}()

	if err != nil {
		log.Printf("Error creating user in Firebase Auth: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user in Firebase Authentication",
		})
	}

	// Map to frontend format and return response
	frontendUser := mappers.MapUserBackendToFrontend(mappers.MapUserFrontendToBackend(&user))
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "You have successfully created an account",
		"user":    frontendUser,
		"token":   token,
	})
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	var loginData models.LoginCredentials

	// Parse request body into LoginCredentials model
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if strings.TrimSpace(loginData.EmailOrUsername) == "" || strings.TrimSpace(loginData.Password) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email/Username and password are required",
		})
	}

	ctx := context.Background()

	var userQuery firestore.Query
	emailOrUsername := strings.TrimSpace(loginData.EmailOrUsername)

	userQuery, err := a.containerService.AuthService.CheckEmailorUsername(ctx, emailOrUsername)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to construct user query",
		})
	}

	// Execute the query
	iter := userQuery.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user data",
		})
	}

	var dbUser map[string]interface{}
	if err := doc.DataTo(&dbUser); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process user data",
		})
	}

	userEmail, ok := dbUser["email"].(string)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user data",
		})
	}

	payload := map[string]string{
		"email":             userEmail,
		"password":          loginData.Password,
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
			"error": "Invalid credentials",
		})
	}
	defer resp.Body.Close()

	// Parse the Firebase response
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var firebaseResponse map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &firebaseResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process authentication response",
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

	// Update user's online status
	backendUser.IsOnline = true
	backendUpdates := mappers.MapUserFrontendToBackend(&backendUser)

	// Update Firestore document
	_, err = doc.Ref.Set(ctx, backendUpdates, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user status",
		})
	}

	// Map backend user to frontend format
	frontendUser := mappers.MapUserToFrontend(&backendUser)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "You have successfully logged in to your account",
		"token":   token,
		"user":    frontendUser,
	})
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {

	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// turn the IsOnline status to false
	ctx := context.Background()
	dbUser, err := a.containerService.AuthService.GetUserByUID(uid)
	dbUser.IsOnline = false

	// Convert the updated User struct to Firestore-compatible format
	backendUpdates := mappers.MapUserFrontendToBackend(dbUser)

	// Update Firestore document
	eror := a.containerService.AuthService.UpdateUser(ctx, uid, backendUpdates)
	if eror != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user status",
		})
	}

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

func (a *AuthHandler) GetSession(c *fiber.Ctx) error {

	// Get token from Authorization header or cookie
	var token string
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		token = c.Cookies("jwt")
		if token != "" {
			log.Printf("Token found in cookie: %s", token[:10])
		} else {
			log.Printf("No token found in either Authorization header or cookie")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "No authentication token provided",
			})
		}
	}

	// Validate token
	uid, err := validateToken(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Path:     "/",
		})
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Token validation failed: %v", err),
		})
	}

	backendUser, err := a.containerService.AuthService.GetUserByUID(uid)
	if err != nil {
		log.Printf("Error fetching user data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user data",
		})
	}

	// Generate new token
	newToken, err := GenerateJWT(backendUser)
	if err != nil {
		log.Printf("Error generating new token: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate token: %v", err),
		})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    newToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	frontendUser := mappers.MapUserToFrontend(backendUser)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user":  frontendUser,
		"token": newToken,
	})
}

// FetchPlatformLogoURL retrieves the platform's logo URL from Firestore
// func FetchPlatformLogoURL() (string, error) {
// 	ctx := context.Background()
// 	logoDoc, err := services.FirestoreClient.Collection("config").Doc("platform").Get(ctx)
// 	if err != nil {
// 		log.Printf("Error fetching platform logo: %v", err)
// 		return "", err
// 	}

// 	// Retrieve the logo URL from the document
// 	logoURL, ok := logoDoc.Data()["logo_url"].(string)
// 	if !ok {
// 		log.Printf("No logo URL found in Firestore")
// 		return "", nil
// 	}

// 	return logoURL, nil
// }
