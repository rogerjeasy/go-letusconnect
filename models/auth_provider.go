package models

// AuthProvider represents the type of authentication provider
type AuthProvider string

const (
	EmailPassword AuthProvider = "email"
	Google        AuthProvider = "google"
	GitHub        AuthProvider = "github"
)

// ProviderData contains provider-specific user information
type ProviderData struct {
	ProviderType AuthProvider   `json:"providerType"`
	Email        string         `json:"email"`
	Password     string         `json:"password,omitempty"`
	DisplayName  string         `json:"displayName"`
	PhotoURL     string         `json:"photoUrl"`
	ProviderID   string         `json:"providerId"`
	Username     string         `json:"username"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	Program      string         `json:"program"`
	ExtraData    map[string]any `json:"extraData,omitempty"`
}
