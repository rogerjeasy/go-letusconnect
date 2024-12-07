package services

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID generates a unique ID
func GenerateID() string {
	id := make([]byte, 16)
	_, _ = rand.Read(id)
	return hex.EncodeToString(id)
}
