package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

func extractProviderData(provider models.AuthProvider, data map[string]interface{}) (models.ProviderData, error) {
	providerData := models.ProviderData{
		ProviderType: provider,
		ExtraData:    make(map[string]any),
	}

	// Helper function to safely extract string from map
	getStringValue := func(m map[string]interface{}, key string) (string, error) {
		if val, ok := m[key]; ok {
			if val == nil {
				return "", fmt.Errorf("%s is required but was nil", key)
			}
			if strVal, ok := val.(string); ok {
				return strVal, nil
			}
			return "", fmt.Errorf("%s must be a string", key)
		}
		return "", fmt.Errorf("%s is required but was missing", key)
	}

	var err error
	switch provider {
	case models.EmailPassword:
		// Handle email/password registration
		if providerData.Email, err = getStringValue(data, "email"); err != nil {
			return providerData, err
		}
		if providerData.Password, err = getStringValue(data, "password"); err != nil {
			return providerData, err
		}
		if providerData.Username, err = getStringValue(data, "username"); err != nil {
			return providerData, err
		}

		// For email/password, use firstName and lastName directly if provided
		firstName, err := getStringValue(data, "firstName")
		if err == nil {
			providerData.FirstName = firstName
		}
		lastName, err := getStringValue(data, "lastName")
		if err == nil {
			providerData.LastName = lastName
		}

		// If names weren't provided, try to extract from username
		if providerData.FirstName == "" && providerData.LastName == "" {
			providerData.FirstName, providerData.LastName = splitDisplayName(providerData.Username)
		}

		if providerData.Program, err = getStringValue(data, "program"); err != nil {
			return providerData, err
		}

	case models.Google:
		// Handle Google OAuth data
		if providerData.Email, err = getStringValue(data, "email"); err != nil {
			return providerData, err
		}

		displayName, err := getStringValue(data, "displayName")
		if err != nil {
			return providerData, err
		}
		providerData.FirstName, providerData.LastName = splitDisplayName(displayName)

		if providerData.PhotoURL, err = getStringValue(data, "photoURL"); err == nil {
			providerData.ExtraData["picture"] = providerData.PhotoURL
		}

		if providerData.ProviderID, err = getStringValue(data, "uid"); err != nil {
			return providerData, err
		}

		providerData.Username = generateUsername(displayName)

		if providerData.Program, err = getStringValue(data, "program"); err != nil {
			return providerData, err
		}

		// Store additional Google-specific data
		if locale, err := getStringValue(data, "locale"); err == nil {
			providerData.ExtraData["locale"] = locale
		}

	case models.GitHub:
		// Handle GitHub OAuth data
		if providerData.Email, err = getStringValue(data, "email"); err != nil {
			return providerData, err
		}

		displayName, err := getStringValue(data, "displayName")
		if err != nil {
			return providerData, err
		}
		providerData.FirstName, providerData.LastName = splitDisplayName(displayName)

		if providerData.PhotoURL, err = getStringValue(data, "photoURL"); err == nil {
			providerData.ExtraData["picture"] = providerData.PhotoURL
		}

		if providerData.ProviderID, err = getStringValue(data, "uid"); err != nil {
			return providerData, err
		}

		if login, err := getStringValue(data, "login"); err == nil {
			providerData.Username = login
			providerData.ExtraData["login"] = login
		} else {
			providerData.Username = generateUsername(displayName)
		}

		if providerData.Program, err = getStringValue(data, "program"); err != nil {
			return providerData, err
		}

		// Store additional GitHub-specific data
		if url, err := getStringValue(data, "html_url"); err == nil {
			providerData.ExtraData["profileUrl"] = url
		}
	}

	return providerData, nil
}

func createFirebaseUser(ctx context.Context, data models.ProviderData) (*auth.UserRecord, error) {
	switch data.ProviderType {
	case models.EmailPassword:
		return services.FirebaseAuth.CreateUser(ctx, (&auth.UserToCreate{}).
			Email(data.Email).
			Password(data.Password))

	case models.Google, models.GitHub:
		return services.FirebaseAuth.GetUserByEmail(ctx, data.Email)
	}

	return nil, fmt.Errorf("unsupported provider type: %s", data.ProviderType)
}

func generateUsername(displayName string) string {
	username := strings.ToLower(strings.ReplaceAll(displayName, " ", "_"))
	return fmt.Sprintf("%s_%d", username, time.Now().UnixNano()%1000)
}

func validateCommonFields(data models.ProviderData) error {
	if strings.TrimSpace(data.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(data.Program) == "" {
		return fmt.Errorf("program is required")
	}
	if data.ProviderType == models.EmailPassword && strings.TrimSpace(data.Password) == "" {
		return fmt.Errorf("password is required for email registration")
	}
	return nil
}

func checkExistingUser(ctx context.Context, data models.ProviderData) error {
	emailQuery := services.Firestore.Collection("users").Where("email", "==", data.Email).Documents(ctx)
	defer emailQuery.Stop()
	if _, err := emailQuery.Next(); err != iterator.Done {
		return fmt.Errorf("email already in use")
	}

	usernameQuery := services.Firestore.Collection("users").Where("username", "==", data.Username).Documents(ctx)
	defer usernameQuery.Stop()
	if _, err := usernameQuery.Next(); err != iterator.Done {
		return fmt.Errorf("username already in use")
	}

	return nil
}

func uploadProfilePicture(ctx context.Context, photoURL string, username string) (string, error) {
	cld := services.CloudinaryClient
	uploadResult, err := cld.Upload.Upload(ctx, photoURL, uploader.UploadParams{
		PublicID: fmt.Sprintf("users/%s/avatar", username),
		Folder:   "users/avatars",
	})
	if err != nil {
		return photoURL, err
	}
	return uploadResult.SecureURL, nil
}

func splitDisplayName(displayName string) (firstName, lastName string) {
	names := strings.Fields(strings.TrimSpace(displayName))

	switch len(names) {
	case 0:
		return "", ""
	case 1:
		return names[0], ""
	default:
		return names[0], strings.Join(names[1:], " ")
	}
}
