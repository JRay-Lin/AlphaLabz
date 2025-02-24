package tools

import (
	"errors"
	"strings"
)

// TokenExtractor extracts the token from a full authorization header string.
//
// It expects the header to be in the format "Bearer <token>".
// If the header is empty or not in the correct format, it returns an error.
func TokenExtractor(fullToken string) (string, error) {
	if fullToken == "" {
		return "", errors.New("missing authorization header")
	}

	// Check if the token starts with "Bearer "
	if strings.HasPrefix(fullToken, "Bearer ") {
		token := strings.TrimSpace(fullToken[7:]) // Trim spaces for safety
		if token == "" {
			return "", errors.New("empty token after Bearer prefix")
		}
		return token, nil
	}

	return "", errors.New("invalid authorization header format")
}
