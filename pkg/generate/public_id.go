package generate

import (
	"crypto/rand"
	"fmt"
)

const publicIDCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// PublicID creates a public ID (e.g., TW-X7Y2P8Q4)
func PublicID() (string, error) {
	const length = 8

	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	publicID := "TW-"
	for i := 0; i < length; i++ {
		publicID += string(publicIDCharset[int(b[i])%len(publicIDCharset)])
	}

	return publicID, nil
}