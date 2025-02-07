package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/config"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Request structures
type LinkedInAuthRequest struct {
	Code string `json:"code"`
}

type LinkedInDataRequest struct {
	LinkedinData struct {
		Profile LinkedInProfileResponse `json:"profile"`
		Email   LinkedInEmailResponse   `json:"email"`
		Token   LinkedInTokenResponse   `json:"token"`
	} `json:"linkedinData"`
}

type LinkedInTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type LinkedInEmailResponse struct {
	Elements []struct {
		Handle struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"handle~"`
	} `json:"elements"`
}

type LinkedInProfileResponse struct {
	ID        string `json:"id"`
	FirstName struct {
		Localized struct {
			EnUS string `json:"en_US"`
		} `json:"localized"`
	} `json:"firstName"`
	LastName struct {
		Localized struct {
			EnUS string `json:"en_US"`
		} `json:"localized"`
	} `json:"lastName"`
}

// LinkedInCallback handles the OAuth callback from LinkedIn
func (a *AuthHandler) LinkedInCallback(c *fiber.Ctx) error {
	var reqData LinkedInDataRequest
	if err := c.BodyParser(&reqData); err != nil {
		// Try parsing as auth request if data request fails
		var authReq LinkedInAuthRequest
		if err := c.BodyParser(&authReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		return a.handleLinkedInAuth(c, authReq)
	}
	return a.handleLinkedInData(c, reqData)
}

// handleLinkedInAuth handles the initial OAuth code exchange
func (a *AuthHandler) handleLinkedInAuth(c *fiber.Ctx, req LinkedInAuthRequest) error {
	// Exchange code for access token
	tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
	tokenData := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		req.Code,
		config.LinkedInRedirectURL,
		config.LinkedInClientID,
		config.LinkedInClientSecret,
	)

	tokenResp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(tokenData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to exchange code for token: %v", err),
		})
	}
	defer tokenResp.Body.Close()

	var tokenResult LinkedInTokenResponse
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenResult); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to parse token response: %v", err),
		})
	}

	// Fetch user data from LinkedIn
	profileData, emailData, err := fetchLinkedInUserData(tokenResult.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch user data: %v", err),
		})
	}

	return a.processLinkedInUser(c, *profileData, *emailData)
}

// handleLinkedInData processes pre-fetched LinkedIn data
func (a *AuthHandler) handleLinkedInData(c *fiber.Ctx, req LinkedInDataRequest) error {
	return a.processLinkedInUser(c, req.LinkedinData.Profile, req.LinkedinData.Email)
}

// fetchLinkedInUserData retrieves user profile and email from LinkedIn
func fetchLinkedInUserData(accessToken string) (*LinkedInProfileResponse, *LinkedInEmailResponse, error) {
	client := &http.Client{}

	// Get user profile
	profileURL := "https://api.linkedin.com/v2/me"
	profileReq, _ := http.NewRequest("GET", profileURL, nil)
	profileReq.Header.Set("Authorization", "Bearer "+accessToken)

	profileResp, err := client.Do(profileReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch profile: %v", err)
	}
	defer profileResp.Body.Close()

	var profileData LinkedInProfileResponse
	if err := json.NewDecoder(profileResp.Body).Decode(&profileData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse profile data: %v", err)
	}

	// Get email address
	emailURL := "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))"
	emailReq, _ := http.NewRequest("GET", emailURL, nil)
	emailReq.Header.Set("Authorization", "Bearer "+accessToken)

	emailResp, err := client.Do(emailReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch email: %v", err)
	}
	defer emailResp.Body.Close()

	var emailData LinkedInEmailResponse
	if err := json.NewDecoder(emailResp.Body).Decode(&emailData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse email data: %v", err)
	}

	return &profileData, &emailData, nil
}

// processLinkedInUser handles the core user creation/login logic
func (a *AuthHandler) processLinkedInUser(c *fiber.Ctx, profileData LinkedInProfileResponse, emailData LinkedInEmailResponse) error {
	// Extract email
	var email string
	if len(emailData.Elements) > 0 {
		email = emailData.Elements[0].Handle.EmailAddress
	}

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No email address found in LinkedIn response",
		})
	}

	// Check if user exists
	existingUser, err := a.authService.GetUserByEmail(email)
	if err != nil {
		// Check if the error is not a "not found" error
		if status, ok := status.FromError(err); !ok || status.Code() != codes.NotFound {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check existing user",
			})
		}
	}

	if existingUser != nil {
		return a.handleExistingUser(c, existingUser)
	}

	return a.createNewLinkedInUser(c, profileData, email)
}

// handleExistingUser processes login for existing users
func (a *AuthHandler) handleExistingUser(c *fiber.Ctx, user *models.User) error {
	// Generate JWT token for existing user
	token, err := GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	setCookie(c, token)

	frontendUser := mappers.MapUserToFrontend(user)
	return c.JSON(fiber.Map{
		"user":    frontendUser,
		"token":   token,
		"message": "Successfully logged in with LinkedIn",
	})
}

func (a *AuthHandler) createNewLinkedInUser(c *fiber.Ctx, profileData LinkedInProfileResponse, email string) error {
	username := strings.Split(email, "@")[0]
	currentTime := time.Now()
	customFormat := "Monday, Jan 2, 2006 at 3:04 PM"

	// Generate profile picture
	profilePictureURL := generateRandomAvatar()
	ctx := context.Background()

	// Upload to Cloudinary
	uploadedURL, err := uploadProfilePicture(ctx, profilePictureURL, username)
	if err != nil {
		log.Printf("Error uploading to Cloudinary: %v", err)
		uploadedURL = profilePictureURL
	}

	// Create user
	newUser := models.User{
		UID:              profileData.ID,
		Username:         username,
		FirstName:        profileData.FirstName.Localized.EnUS,
		LastName:         profileData.LastName.Localized.EnUS,
		Email:            email,
		ProfilePicture:   uploadedURL,
		AccountCreatedAt: FormatTime(currentTime, customFormat),
		IsActive:         true,
		IsVerified:       true,
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

	// Save to Firestore
	if err := a.authService.CreateUser(ctx, &newUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save user",
		})
	}

	// Generate token
	token, err := GenerateJWT(&newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	setCookie(c, token)

	// Send welcome email
	go func() {
		if err := SendWelcomeEmail(newUser.Email, newUser.Username, "linkedin"); err != nil {
			log.Printf("Error sending welcome email: %v", err)
		}
	}()

	// Send notification
	go func() {
		if err := a.containerService.GeneralNotificationService.SendNewUserNotification(context.Background(), &newUser); err != nil {
			log.Printf("Failed to send new user notification: %v", err)
		}
	}()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully created account with LinkedIn",
		"user":    mappers.MapUserToFrontend(&newUser),
		"token":   token,
	})
}

// Helper function to set JWT cookie
func setCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: "Lax",
	})
}
